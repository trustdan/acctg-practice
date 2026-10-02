package tui

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"github.com/trustdan/acctg-practice/internal/storage"
)

// ShipRoll represents the 3D roll angle of the spaceship.
type ShipRoll int

const (
	RollHardDown ShipRoll = -2 // Hard bank down: ventral armor plates visible in 3D perspective
	RollSoftDown ShipRoll = -1 // Gentle bank down: lower wing tilted into view
	RollLevel    ShipRoll = 0  // Level flight: razor-sleek neutral side profile
	RollSoftUp   ShipRoll = 1  // Gentle bank up: upper deck and canopy tilt into view
	RollHardUp   ShipRoll = 2  // Hard bank up: dorsal delta wing and dual cockpit glass in 3D perspective
)

const (
	// FireHoldTicks (~600ms) outlasts the OS key-repeat delay. Terminals send no
	// key-release events, and pressing a steering key while F is held stops F
	// from repeating, so steering keys refresh this window instead.
	FireHoldTicks = 18

	// A fiscal year is four 40-second quarters; closing one wins the round.
	QuarterTicks    = 1200
	QuartersPerYear = 4
	YearTicks       = QuarterTicks * QuartersPerYear

	carePackageEvery = 540 // ~18s between field audit kit drops
	maxBombs         = 5
)

// introTickMsg drives the ~30 FPS startup animation loop.
type introTickMsg struct{}

func introTick() tea.Cmd {
	return tea.Tick(33*time.Millisecond, func(t time.Time) tea.Msg {
		return introTickMsg{}
	})
}

// IntroStar represents a background star drifting in parallax.
type IntroStar struct {
	X     float64
	Y     float64
	Speed float64
	Layer int // 0: dim distant, 1: medium cyan, 2: bright warp streak
}

// IntroLaser represents a high-energy blaster projectile fired from the ship.
type IntroLaser struct {
	X float64
	Y float64
}

// IntroShockwave represents an expanding smart-bomb detonation ring.
type IntroShockwave struct {
	X         float64
	Y         float64
	Radius    float64
	MaxRadius float64
	Life      int
}

// IntroTargetType categorizes incoming accounting objects.
type IntroTargetType int

const (
	TargetTAccount IntroTargetType = iota
	TargetDelta
	TargetSum
	TargetConcept
	TargetAsteroid
	TargetFragment
	TargetCarePackage
)

// PickupKind identifies a care package the ship collects by flying into it.
type PickupKind int

const (
	PickupNone   PickupKind = iota
	PickupRepair            // restores shield and external audit integrity
	PickupBomb              // restocks one smart bomb
)

// IntroTarget represents an incoming accounting entity or hazard drifting from right to left.
type IntroTarget struct {
	ID         int
	Type       IntroTargetType
	Text       string
	X          float64
	Y          float64
	Speed      float64
	Width      int
	Color      lipgloss.Color
	MaxHP      int
	HP         int
	IsHazard   bool
	HazardName string
	Pickup     PickupKind
}

// IntroParticle represents debris from an exploded accounting entity or asteroid.
type IntroParticle struct {
	X     float64
	Y     float64
	VX    float64
	VY    float64
	Life  int
	Char  rune
	Color lipgloss.Color
}

// IntroCallout represents floating combat / auditing text popping up on impact.
type IntroCallout struct {
	Text  string
	X     float64
	Y     float64
	Life  int
	Color lipgloss.Color
}

// IntroState manages the complete spaceship flight, targeting, combat, and audit health simulation.
type IntroState struct {
	Width        int
	Height       int
	RNG          *rand.Rand
	TickCount    int
	Score        int
	BlastedCount int
	NextTargetID int

	// Spaceship flight dynamics
	ShipX           float64
	ShipY           float64
	ShipVY          float64
	ShipRoll        ShipRoll
	TargetY         float64
	LastManual      time.Time
	ManualMode      bool
	LaserCool       int
	RapidFireRounds int  // ticks the trigger is still considered held
	AutoFire        bool // trigger lock: fire continuously until toggled off
	BombsRemaining  int
	Shockwaves      []IntroShockwave
	Paused          bool

	// Progressive speed acceleration & difficulty
	SpeedFactor float64

	// Fiscal-year campaign: survive four quarters to close the books.
	Year         int // also the score multiplier
	YearTick     int
	YearComplete bool // year-end "continue?" prompt is showing
	YearBonus    int
	Victory      bool // run ended by retiring after a closed year

	// Dual Audit Health System
	ShieldHP       int // Internal Audit Shields (100 -> 0)
	GlobalHP       int // External Audit Integrity (100 -> 0)
	GameOver       bool
	GameOverReason string
	lastHitTick    int

	// High Score & Retro 3-Initials Hall of Fame
	HighScore           int
	HighScoreHolder     string
	InitialsEntryActive bool
	Initials            [3]rune
	InitialsCursor      int
	ScoreSaved          bool

	// Hall of Fame overlay; flight is paused while it is open.
	ShowLeaderboard bool
	LeaderboardOnly bool // opened from outside the game, so closing returns to the previous screen
	Leaderboard     []storage.ArcadeHighScore

	// Simulation collections
	Stars     []IntroStar
	Lasers    []IntroLaser
	Targets   []IntroTarget
	Particles []IntroParticle
	Callouts  []IntroCallout
}

// NewIntroState creates and initializes a new startup animation simulation.
func NewIntroState(width, height int, rng *rand.Rand) *IntroState {
	if width < 60 {
		width = 80
	}
	if height < 18 {
		height = 24
	}
	if rng == nil {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	playTop := 3
	playBottom := height - 4
	midY := float64(playTop+playBottom) / 2.0

	s := &IntroState{
		Width:           width,
		Height:          height,
		RNG:             rng,
		ShipX:           4.0,
		ShipY:           midY,
		ShipVY:          0.0,
		ShipRoll:        RollLevel,
		TargetY:         midY,
		TickCount:       0,
		BombsRemaining:  3,
		ShieldHP:        100,
		GlobalHP:        100,
		SpeedFactor:     1.0,
		Year:            1,
		HighScore:       5000,
		HighScoreHolder: "DAN",
		Initials:        [3]rune{'A', 'A', 'A'},
	}

	s.initStars()
	s.spawnInitialTargets()
	return s
}

func (s *IntroState) initStars() {
	starCount := (s.Width * s.Height) / 25
	if starCount < 30 {
		starCount = 30
	}
	s.Stars = make([]IntroStar, starCount)
	for i := 0; i < starCount; i++ {
		layer := s.RNG.Intn(3)
		speed := 0.18
		if layer == 1 {
			speed = 0.42
		} else if layer == 2 {
			speed = 0.85
		}
		s.Stars[i] = IntroStar{
			X:     float64(s.RNG.Intn(s.Width)),
			Y:     float64(3 + s.RNG.Intn(max(1, s.Height-6))),
			Speed: speed,
			Layer: layer,
		}
	}
}

func (s *IntroState) spawnInitialTargets() {
	for i := 0; i < 3; i++ {
		x := float64(s.Width - 15 - i*16)
		s.spawnTargetAt(x)
	}
}

// Resize dynamically recalculates screen boundaries and starfield distribution.
func (s *IntroState) Resize(width, height int) {
	if width < 60 {
		width = 80
	}
	if height < 18 {
		height = 24
	}
	s.Width = width
	s.Height = height

	playTop := 5.0
	playBottom := float64(s.Height - 5)
	if s.ShipY < playTop {
		s.ShipY = playTop
	}
	if s.ShipY > playBottom {
		s.ShipY = playBottom
	}

	s.initStars()
}

// inputLocked reports whether flight controls are ignored right now.
func (s *IntroState) inputLocked() bool {
	return s.GameOver || s.InitialsEntryActive || s.Paused || s.YearComplete
}

// SteerUp moves the spaceship upward manually.
func (s *IntroState) SteerUp() {
	s.ManualMode = true
	if s.inputLocked() {
		return
	}
	s.ShipVY = -0.75
	s.ShipRoll = RollHardUp
	s.LastManual = time.Now()
	s.holdTrigger()
}

// SteerDown moves the spaceship downward manually.
func (s *IntroState) SteerDown() {
	s.ManualMode = true
	if s.inputLocked() {
		return
	}
	s.ShipVY = 0.75
	s.ShipRoll = RollHardDown
	s.LastManual = time.Now()
	s.holdTrigger()
}

// holdTrigger keeps an active machine-gun burst alive while steering keys
// repeat, since the terminal stops repeating F once another key goes down.
func (s *IntroState) holdTrigger() {
	if s.RapidFireRounds > 0 {
		s.RapidFireRounds = FireHoldTicks
	}
}

// ToggleAutoFire locks the trigger on (or off) for hands-free machine-gun fire.
func (s *IntroState) ToggleAutoFire() {
	s.ManualMode = true
	if s.inputLocked() {
		return
	}
	s.AutoFire = !s.AutoFire
	msg := "AUTO-FIRE OFF"
	if s.AutoFire {
		msg = "AUTO-FIRE LOCKED ON [G]"
	}
	s.spawnCallout(s.ShipX+10.0, s.ShipY-2.0, msg, lipgloss.Color("#00FFFF"))
}

// TogglePause freezes or resumes the flight.
func (s *IntroState) TogglePause() {
	if s.GameOver || s.InitialsEntryActive || s.YearComplete {
		return
	}
	s.Paused = !s.Paused
}

// Quarter returns the current fiscal quarter (1-4).
func (s *IntroState) Quarter() int {
	return min(QuartersPerYear, s.YearTick/QuarterTicks+1)
}

// QuarterSecondsLeft returns the whole seconds remaining in the current quarter.
func (s *IntroState) QuarterSecondsLeft() int {
	left := QuarterTicks - s.YearTick%QuarterTicks
	if s.YearTick >= YearTicks {
		left = 0
	}
	return (left + 29) / 30
}

// ContinueYear answers "yes" at the year-end prompt: a faster year with a
// bigger score multiplier, after a partial repair and bomb restock.
func (s *IntroState) ContinueYear() {
	if !s.YearComplete {
		return
	}
	s.Year++
	s.YearTick = 0
	s.YearComplete = false
	s.ShieldHP = min(100, s.ShieldHP+40)
	s.GlobalHP = min(100, s.GlobalHP+40)
	s.BombsRemaining = min(maxBombs, s.BombsRemaining+2)
	s.spawnCallout(s.ShipX+10.0, s.ShipY-2.0, fmt.Sprintf("FISCAL YEAR %d // SCORE x%d", s.Year, s.Year), lipgloss.Color("#FFE600"))
}

// RetireAfterYear answers "no" at the year-end prompt and ends the run as a win.
func (s *IntroState) RetireAfterYear() {
	if !s.YearComplete {
		return
	}
	s.YearComplete = false
	s.Victory = true
	s.GameOver = true
	s.GameOverReason = fmt.Sprintf("RETIRED AFTER %d FISCAL YEAR(S) WITH A CLEAN OPINION", s.Year)
}

// addScore applies the fiscal-year multiplier.
func (s *IntroState) addScore(points int) {
	s.Score += points * max(1, s.Year)
}

// Fire triggers manual blaster cannons and sets rapid-fire burst counter for continuous machine-gun bursts.
func (s *IntroState) Fire() {
	s.ManualMode = true
	if s.inputLocked() {
		return
	}
	s.RapidFireRounds = FireHoldTicks
	if s.LaserCool <= 0 {
		s.Lasers = append(s.Lasers,
			IntroLaser{X: s.ShipX + 19.0, Y: s.ShipY - 0.5},
			IntroLaser{X: s.ShipX + 19.0, Y: s.ShipY + 0.5},
		)
		s.LaserCool = 2
	}
}

// DeployBomb detonates a smart bomb, launching an expanding shockwave that vaporizes standard threats
// and heavily damages large accounting asteroids.
func (s *IntroState) DeployBomb() bool {
	s.ManualMode = true
	if s.inputLocked() {
		return false
	}
	if s.BombsRemaining <= 0 {
		s.spawnCallout(s.ShipX+10.0, s.ShipY, "NO BOMBS REMAINING! [0 LEFT]", lipgloss.Color("#EF4444"))
		return false
	}
	s.BombsRemaining--

	// 1. Expanding shockwave ring
	s.Shockwaves = append(s.Shockwaves, IntroShockwave{
		X:         s.ShipX + 16.0,
		Y:         s.ShipY,
		Radius:    2.0,
		MaxRadius: float64(max(s.Width, 80)),
		Life:      12,
	})

	// 2. Radial explosive debris
	chars := []rune{'💥', '░', '▒', '▓', '*', '+', 'o', '·'}
	colors := []lipgloss.Color{
		lipgloss.Color("#FFFFFF"),
		lipgloss.Color("#FFE600"),
		lipgloss.Color("#FF5500"),
		lipgloss.Color("#00FFFF"),
		lipgloss.Color("#A855F7"),
	}
	for i := 0; i < 32; i++ {
		ang := s.RNG.Float64() * 2 * math.Pi
		spd := 0.5 + s.RNG.Float64()*1.8
		s.Particles = append(s.Particles, IntroParticle{
			X:     s.ShipX + 16.0,
			Y:     s.ShipY,
			VX:    math.Cos(ang) * spd * 1.8,
			VY:    math.Sin(ang) * spd * 0.9,
			Life:  10 + s.RNG.Intn(8),
			Char:  chars[s.RNG.Intn(len(chars))],
			Color: colors[s.RNG.Intn(len(colors))],
		})
	}

	// 3. Vaporize or heavily damage all active targets on screen
	var remainingTargets []IntroTarget
	for _, t := range s.Targets {
		if t.Pickup != PickupNone {
			remainingTargets = append(remainingTargets, t)
			continue
		}
		if t.IsHazard {
			t.HP -= 15
			if t.HP <= 0 {
				s.spawnExplosion(t.X+float64(t.Width)/2.0, t.Y)
				frags := s.spawnFragments(t.X, t.Y, t.HazardName)
				remainingTargets = append(remainingTargets, frags...)
				s.addScore(1000)
				s.BlastedCount++
				s.spawnCallout(t.X, t.Y, "HAZARD VAPORIZED! +1,000", lipgloss.Color("#FFE600"))
			} else {
				s.spawnExplosion(t.X+float64(t.Width)/2.0, t.Y)
				s.spawnCallout(t.X, t.Y, fmt.Sprintf("-15 HP! [%d/%d]", t.HP, t.MaxHP), lipgloss.Color("#EF4444"))
				remainingTargets = append(remainingTargets, t)
			}
		} else {
			s.spawnExplosion(t.X+float64(t.Width)/2.0, t.Y)
			s.addScore(100)
			s.BlastedCount++
		}
	}
	s.Targets = remainingTargets

	s.spawnCallout(s.ShipX+18.0, s.ShipY-1.0, "🚨 SMART BOMB DETONATED! TOTAL AUDIT CLEARANCE! 🚨", lipgloss.Color("#00FFFF"))
	return true
}

func (s *IntroState) spawnFragments(x, y float64, hazardName string) []IntroTarget {
	fragDefs := []struct {
		text  string
		score int
		color lipgloss.Color
	}{
		{"[SHREDDED EVIDENCE]", 200, lipgloss.Color("#38BDF8")},
		{"[SEC PENALTY]", 300, lipgloss.Color("#F59E0B")},
		{"[RESTITUTION]", 250, lipgloss.Color("#10B981")},
	}
	var frags []IntroTarget
	for i, fd := range fragDefs {
		yOffset := float64(i-1) * 1.5
		s.NextTargetID++
		frags = append(frags, IntroTarget{
			ID:    s.NextTargetID,
			Type:  TargetFragment,
			Text:  fd.text,
			X:     x + float64(i*8),
			Y:     math.Max(4.0, math.Min(float64(s.Height-5), y+yOffset)),
			Speed: 0.45 * s.SpeedFactor,
			Width: textWidth(fd.text),
			Color: fd.color,
			MaxHP: 1,
			HP:    1,
		})
	}
	return frags
}

// SetManualMode explicitly sets or toggles manual piloting.
func (s *IntroState) SetManualMode(manual bool) {
	s.ManualMode = manual
}

// SetHighScore updates the current all-time high score and holder.
func (s *IntroState) SetHighScore(score int, holder string) {
	s.HighScore = score
	if holder != "" {
		s.HighScoreHolder = holder
	}
}

// CycleInitial increments or decrements the letter at the current initials cursor.
func (s *IntroState) CycleInitial(direction int) {
	if s.InitialsCursor < 0 || s.InitialsCursor > 2 {
		s.InitialsCursor = 0
	}
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	cur := s.Initials[s.InitialsCursor]
	idx := strings.IndexRune(chars, cur)
	if idx == -1 {
		idx = 0
	}
	if direction > 0 {
		idx = (idx + 1) % len(chars)
	} else {
		idx = (idx - 1 + len(chars)) % len(chars)
	}
	s.Initials[s.InitialsCursor] = rune(chars[idx])
}

// SetInitial sets a character directly and advances the cursor.
func (s *IntroState) SetInitial(ch rune) {
	if s.InitialsCursor < 0 || s.InitialsCursor > 2 {
		s.InitialsCursor = 0
	}
	upper := strings.ToUpper(string(ch))
	if len(upper) > 0 {
		r := rune(upper[0])
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			s.Initials[s.InitialsCursor] = r
			if s.InitialsCursor < 2 {
				s.InitialsCursor++
			}
		}
	}
}

// PrevInitial moves the initials cursor back one space.
func (s *IntroState) PrevInitial() {
	if s.InitialsCursor > 0 {
		s.InitialsCursor--
	}
}

// NextInitial advances the initials cursor, returning true if submitting from the final slot.
func (s *IntroState) NextInitial() bool {
	if s.InitialsCursor < 2 {
		s.InitialsCursor++
		return false
	}
	return true
}

// GetInitialsString returns the 3-letter initials string.
func (s *IntroState) GetInitialsString() string {
	return string(s.Initials[:])
}

func (s *IntroState) triggerGameOver(reason string) {
	if s.GameOver {
		return
	}
	s.GameOver = true
	s.GameOverReason = reason
	s.spawnExplosion(s.ShipX+9.0, s.ShipY)
	s.spawnExplosion(s.ShipX+4.0, s.ShipY-1.0)
	s.spawnExplosion(s.ShipX+14.0, s.ShipY+1.0)
	s.spawnCallout(s.ShipX+5.0, s.ShipY-2.0, "💥 SHIP COMPROMISED!", lipgloss.Color("#EF4444"))
}

// Update advances the simulation state by one tick (~33ms).
func (s *IntroState) Update() {
	if s.ShowLeaderboard || s.Paused {
		return
	}
	s.TickCount++
	playTop := 5.0
	playBottom := float64(s.Height - 5)
	midY := (playTop + playBottom) / 2.0

	// The fiscal-year clock only runs once a player has taken the controls.
	if s.ManualMode && !s.GameOver && !s.YearComplete {
		prevQuarter := s.Quarter()
		s.YearTick++
		if s.YearTick >= YearTicks {
			s.closeYear()
		} else if q := s.Quarter(); q != prevQuarter {
			s.spawnCallout(s.ShipX+10.0, s.ShipY-2.0, fmt.Sprintf("Q%d BEGINS // AUDIT KIT INBOUND", q), lipgloss.Color("#FFE600"))
			s.spawnCarePackage(PickupRepair)
		} else if s.YearTick%carePackageEvery == 0 {
			kind := PickupRepair
			if s.RNG.Intn(4) == 0 {
				kind = PickupBomb
			}
			s.spawnCarePackage(kind)
		}
		s.regenerate()
	}
	active := !s.GameOver && !s.YearComplete

	// Gentle acceleration across the year; each continued year starts faster.
	s.SpeedFactor = math.Min(3.0, 1.0+0.9*float64(s.YearTick)/YearTicks+0.4*float64(s.Year-1))

	// Update Shockwaves
	activeShockwaves := s.Shockwaves[:0]
	for _, sw := range s.Shockwaves {
		sw.Radius += 4.5
		sw.Life--
		if sw.Life > 0 && sw.Radius < sw.MaxRadius {
			activeShockwaves = append(activeShockwaves, sw)
		}
	}
	s.Shockwaves = activeShockwaves

	// 1. Spaceship Physics & 3D Roll calculation
	if active {
		if !s.ManualMode {
			// Auto-pilot (Demo / Attract mode): search for incoming target ahead
			var bestTarget *IntroTarget
			bestDist := 9999.0
			for i := range s.Targets {
				t := &s.Targets[i]
				if t.X > s.ShipX+10.0 {
					dist := t.X - s.ShipX
					if dist < bestDist {
						bestDist = dist
						bestTarget = t
					}
				}
			}

			if bestTarget != nil {
				s.TargetY = bestTarget.Y
			} else {
				sineWave := math.Sin(float64(s.TickCount)*0.07) * (playBottom - playTop) * 0.36
				cosineSway := math.Cos(float64(s.TickCount)*0.03) * 1.5
				s.TargetY = midY + sineWave + cosineSway
			}

			diff := s.TargetY - s.ShipY
			desiredVY := diff * 0.16
			if desiredVY > 0.65 {
				desiredVY = 0.65
			}
			if desiredVY < -0.65 {
				desiredVY = -0.65
			}
			s.ShipVY += (desiredVY - s.ShipVY) * 0.25
		} else {
			// Manual flight mode: gentle aerodynamic drag settles ship in lane when keys released
			s.ShipVY *= 0.88
			if math.Abs(s.ShipVY) < 0.03 {
				s.ShipVY = 0
			}
		}

		// Apply vertical velocity with boundary damping
		s.ShipY += s.ShipVY
		if s.ShipY < playTop {
			s.ShipY = playTop
			s.ShipVY = 0.1
		}
		if s.ShipY > playBottom {
			s.ShipY = playBottom
			s.ShipVY = -0.1
		}

		// Calculate 3D Roll angle based on vertical velocity
		if s.ShipVY < -0.42 {
			s.ShipRoll = RollHardUp
		} else if s.ShipVY < -0.10 {
			s.ShipRoll = RollSoftUp
		} else if s.ShipVY > 0.42 {
			s.ShipRoll = RollHardDown
		} else if s.ShipVY > 0.10 {
			s.ShipRoll = RollSoftDown
		} else {
			s.ShipRoll = RollLevel
		}

		// 2. Spaceship Blasters & Rapid-Fire Machine Gun
		if !s.ManualMode {
			// Auto-fire only in demo mode
			s.LaserCool--
			if s.LaserCool <= 0 {
				aligned := false
				for _, t := range s.Targets {
					if t.X > s.ShipX+10.0 && t.X < float64(s.Width) && math.Abs(t.Y-s.ShipY) <= 1.8 {
						aligned = true
						break
					}
				}
				if aligned || s.TickCount%8 == 0 {
					s.Lasers = append(s.Lasers, IntroLaser{X: s.ShipX + 27.0, Y: s.ShipY})
					s.LaserCool = 6
				}
			}
		} else {
			// In manual mode, the machine gun fires while the trigger is held or locked on
			s.LaserCool--
			if s.RapidFireRounds > 0 || s.AutoFire {
				if s.RapidFireRounds > 0 {
					s.RapidFireRounds--
				}
				if s.LaserCool <= 0 {
					s.Lasers = append(s.Lasers,
						IntroLaser{X: s.ShipX + 27.0, Y: s.ShipY - 0.8},
						IntroLaser{X: s.ShipX + 27.0, Y: s.ShipY + 0.8},
					)
					s.LaserCool = 2
				}
			}
		}
	}

	// 3. Update Lasers
	activeLasers := s.Lasers[:0]
	for _, l := range s.Lasers {
		l.X += 2.4
		if l.X < float64(s.Width) {
			activeLasers = append(activeLasers, l)
		}
	}
	s.Lasers = activeLasers

	// 4. Update Parallax Starfield (speed scales with acceleration)
	for i := range s.Stars {
		s.Stars[i].X -= s.Stars[i].Speed * s.SpeedFactor
		if s.Stars[i].X < 0 {
			s.Stars[i].X = float64(s.Width - 1)
			s.Stars[i].Y = float64(3 + s.RNG.Intn(max(1, s.Height-6)))
		}
	}

	// 5. Spawn Incoming Accounting Targets and Asteroid Hazards
	if active {
		targetCap := 5
		if s.Width > 90 {
			targetCap = 7
		}

		spawnInterval := max(5, int(14.0/s.SpeedFactor))
		if len(s.Targets) < targetCap && (s.TickCount%spawnInterval == 0 || len(s.Targets) == 0) {
			s.spawnTargetAt(float64(s.Width - 2))
		}

		// Spawn heavy asteroid hazards periodically
		hazardInterval := max(45, int(95.0/s.SpeedFactor))
		if s.TickCount%hazardInterval == 0 {
			s.spawnHazardAsteroid()
		}
	}

	// 6. Update Targets, Laser Collisions, Ship Collisions, and Off-Screen Escapes
	activeTargets := s.Targets[:0]
	for _, t := range s.Targets {
		t.X -= t.Speed * s.SpeedFactor

		// A. Check Collision with Ship Fuselage (Internal Audit Shields)
		if active {
			shipMinX := s.ShipX
			shipY := s.ShipY
			// Only the fuselage row reaches the nose; the swept wings end further back.
			shipMaxX := s.ShipX + 21.0
			if math.Abs(t.Y-shipY) < 0.5 {
				shipMaxX = s.ShipX + 27.0
			}

			if math.Abs(t.Y-shipY) <= 1.5 && t.X <= shipMaxX && t.X+float64(t.Width) >= shipMinX {
				if t.Pickup != PickupNone {
					s.collect(t.Pickup)
					continue
				}
				s.lastHitTick = s.TickCount
				dmg := 20
				if t.IsHazard {
					dmg = 35
				}
				s.ShieldHP -= dmg
				s.spawnExplosion(t.X+float64(t.Width)/2.0, t.Y)
				if t.IsHazard {
					s.spawnCallout(s.ShipX+10.0, s.ShipY-1.0, fmt.Sprintf("⚠ HAZARD COLLISION! -%d%% SHIELD", dmg), lipgloss.Color("#EF4444"))
				} else {
					s.spawnCallout(s.ShipX+10.0, s.ShipY-1.0, fmt.Sprintf("SHIELD HIT! -%d%%", dmg), lipgloss.Color("#F59E0B"))
				}

				if s.ShieldHP <= 0 {
					s.ShieldHP = 0
					s.triggerGameOver("INTERNAL AUDIT FAILURE: MATERIAL DEFICIENCY IN CONTROLS")
				}
				continue // target destroyed by impact
			}
		}

		// B. Check Laser Collisions (care packages are not shootable)
		hit := false
		hitLaserIdx := -1
		for lIdx, l := range s.Lasers {
			if t.Pickup != PickupNone {
				break
			}
			if math.Abs(l.Y-t.Y) <= 1.2 && l.X >= t.X && l.X <= t.X+float64(t.Width)+1.5 {
				hit = true
				hitLaserIdx = lIdx
				break
			}
		}

		if hit {
			if hitLaserIdx >= 0 && hitLaserIdx < len(s.Lasers) {
				s.Lasers = append(s.Lasers[:hitLaserIdx], s.Lasers[hitLaserIdx+1:]...)
			}

			if t.IsHazard {
				t.HP--
				s.spawnSpark(t.X+float64(t.Width)/2.0, t.Y)
				s.spawnCallout(t.X+float64(t.Width)/2.0, t.Y, fmt.Sprintf("-1 HP [%d/%d]", t.HP, t.MaxHP), lipgloss.Color("#F59E0B"))

				if t.HP <= 0 {
					s.spawnExplosion(t.X+float64(t.Width)/2.0, t.Y)
					frags := s.spawnFragments(t.X, t.Y, t.HazardName)
					activeTargets = append(activeTargets, frags...)
					s.addScore(1000)
					s.BlastedCount++
					s.spawnCallout(t.X, t.Y, "🚨 HAZARD NEUTRALIZED! +1,000 🚨", lipgloss.Color("#FFE600"))
					continue
				} else {
					t.Text = fmt.Sprintf("[%s: HP %d/%d]", t.HazardName, t.HP, t.MaxHP)
					t.Width = textWidth(t.Text)
					activeTargets = append(activeTargets, t)
					continue
				}
			} else {
				s.spawnExplosion(t.X+float64(t.Width)/2.0, t.Y)
				s.spawnCallout(t.X, t.Y)
				if t.Type == TargetFragment {
					s.addScore(250)
				} else {
					s.addScore(100)
				}
				s.BlastedCount++
				continue
			}
		}

		// C. Check Escaping past left edge un-audited (External Audit Global HP)
		if t.X+float64(t.Width) <= 0 {
			if active && t.Pickup == PickupNone {
				dmg := 10
				if t.IsHazard {
					dmg = 25
					s.spawnCallout(2.0, t.Y, "🚨 AUDIT BREACH! HAZARD ESCAPED -25%", lipgloss.Color("#EF4444"))
				} else {
					s.spawnCallout(2.0, t.Y, "UNAUDITED SLIP! -10% EXT AUDIT", lipgloss.Color("#F59E0B"))
				}
				s.GlobalHP -= dmg
				if s.GlobalHP <= 0 {
					s.GlobalHP = 0
					s.triggerGameOver("EXTERNAL AUDIT FAILURE: ADVERSE OPINION (UNAUDITED ENTITIES)")
				}
			}
			continue
		}

		activeTargets = append(activeTargets, t)
	}
	s.Targets = activeTargets

	// 7. Update Particles
	activeParticles := s.Particles[:0]
	for _, p := range s.Particles {
		p.X += p.VX
		p.Y += p.VY
		p.Life--
		if p.Life > 0 && p.X >= 0 && p.X < float64(s.Width) && p.Y >= 2 && p.Y < float64(s.Height-1) {
			activeParticles = append(activeParticles, p)
		}
	}
	s.Particles = activeParticles

	// 8. Update Floating Callouts
	activeCallouts := s.Callouts[:0]
	for _, c := range s.Callouts {
		c.Y -= 0.18
		c.Life--
		if c.Life > 0 && c.Y >= 2 {
			activeCallouts = append(activeCallouts, c)
		}
	}
	s.Callouts = activeCallouts
}

// regenerate slowly restores shields after a few seconds without a hit, and
// external audit integrity at a slower steady rate.
func (s *IntroState) regenerate() {
	if s.TickCount-s.lastHitTick >= 90 && s.TickCount%45 == 0 {
		s.ShieldHP = min(100, s.ShieldHP+1)
	}
	if s.TickCount%150 == 0 {
		s.GlobalHP = min(100, s.GlobalHP+1)
	}
}

// collect applies a care package the ship flew into.
func (s *IntroState) collect(kind PickupKind) {
	switch kind {
	case PickupRepair:
		s.ShieldHP = min(100, s.ShieldHP+30)
		s.GlobalHP = min(100, s.GlobalHP+15)
		s.spawnCallout(s.ShipX+10.0, s.ShipY-2.0, "✚ AUDIT KIT: +30 SHIELD / +15 EXT AUDIT", lipgloss.Color("#22C55E"))
	case PickupBomb:
		s.BombsRemaining = min(maxBombs, s.BombsRemaining+1)
		s.spawnCallout(s.ShipX+10.0, s.ShipY-2.0, "✚ SMART BOMB RESTOCKED!", lipgloss.Color("#00FFFF"))
	}
	s.addScore(150)
}

// closeYear ends the fiscal year: clears the field and banks a bonus for the
// health left on both audit meters.
func (s *IntroState) closeYear() {
	s.YearComplete = true
	s.YearBonus = (5000 + (s.ShieldHP+s.GlobalHP)*25) * s.Year
	s.Score += s.YearBonus
	for _, t := range s.Targets {
		s.spawnExplosion(t.X+float64(t.Width)/2.0, t.Y)
	}
	s.Targets = nil
	s.Lasers = nil
	s.RapidFireRounds = 0
}

func (s *IntroState) spawnCarePackage(kind PickupKind) {
	s.NextTargetID++
	yRange := max(1, s.Height-10)
	text, color := "[✚ AUDIT KIT ✚]", lipgloss.Color("#22C55E")
	if kind == PickupBomb {
		text, color = "[✚ +1 BOMB ✚]", lipgloss.Color("#00FFFF")
	}
	s.Targets = append(s.Targets, IntroTarget{
		ID:     s.NextTargetID,
		Type:   TargetCarePackage,
		Text:   text,
		X:      float64(s.Width - 2),
		Y:      float64(5 + s.RNG.Intn(yRange)),
		Speed:  0.26,
		Width:  textWidth(text),
		Color:  color,
		MaxHP:  1,
		HP:     1,
		Pickup: kind,
	})
}

func (s *IntroState) spawnTargetAt(x float64) {
	s.NextTargetID++
	playTop := 4
	playBottom := s.Height - 4
	yRange := playBottom - playTop
	if yRange < 1 {
		yRange = 1
	}
	y := float64(playTop + s.RNG.Intn(yRange))

	type targetDef struct {
		targetType IntroTargetType
		text       string
		color      lipgloss.Color
	}

	amber := lipgloss.Color("#F59E0B")
	skyCyan := lipgloss.Color("#38BDF8")
	emerald := lipgloss.Color("#10B981")
	magenta := lipgloss.Color("#F43F5E")
	purple := lipgloss.Color("#A855F7")
	lime := lipgloss.Color("#84CC16")
	gold := lipgloss.Color("#EAB308")

	catalog := []targetDef{
		{TargetTAccount, "[─┬─ CASH ─┬─]", emerald},
		{TargetTAccount, "[ Dr│ACCTS REC│Cr ]", skyCyan},
		{TargetTAccount, "[─┬─ EQUIPMENT ─┬─]", amber},
		{TargetTAccount, "[ Dr│INVENTORY│Cr ]", purple},
		{TargetTAccount, "[─┬─ PREPAID RENT ─┬─]", skyCyan},
		{TargetTAccount, "[ Cr│ACCTS PAYABLE│Dr ]", amber},
		{TargetTAccount, "[─┬─ COMMON STOCK ─┬─]", skyCyan},
		{TargetTAccount, "[ Cr│SERVICE REVENUE│Dr ]", emerald},
		{TargetTAccount, "[ Dr│DIVIDENDS│Cr ]", magenta},
		{TargetTAccount, "[─┬─ UNEARNED REV ─┬─]", purple},
		{TargetTAccount, "[ Dr│WAGE EXPENSE│Cr ]", magenta},
		{TargetDelta, "ΔA = ΔL + ΔE", magenta},
		{TargetDelta, "Δ CASH = +$5,000", emerald},
		{TargetDelta, "Δ ASSETS", skyCyan},
		{TargetDelta, "Δ LIABILITIES", amber},
		{TargetDelta, "Δ RETAINED EARNINGS", purple},
		{TargetDelta, "Δ EQUITY", skyCyan},
		{TargetDelta, "Δ WORKING CAPITAL", lime},
		{TargetSum, "Σ Dr = Σ Cr", emerald},
		{TargetSum, "Σ ASSETS = $100K", skyCyan},
		{TargetSum, "Σ DEBITS", emerald},
		{TargetSum, "Σ CREDITS", amber},
		{TargetSum, "+$25,000 CASH", lime},
		{TargetSum, "-$12,000 PAYABLE", magenta},
		{TargetSum, "$$$ $9,000,000 $$$", gold},
		{TargetConcept, "<GAAP COMPLIANCE>", emerald},
		{TargetConcept, "<UNEARNED REVENUE>", purple},
		{TargetConcept, "<ACCRUED EXPENSE>", amber},
		{TargetConcept, "<DEPRECIATION>", magenta},
		{TargetConcept, "<TRIAL BALANCE>", skyCyan},
		{TargetConcept, "<AUDIT PROOF>", emerald},
		{TargetConcept, "<DOUBLE ENTRY>", skyCyan},
	}

	chosen := catalog[s.RNG.Intn(len(catalog))]
	speed := 0.32 + s.RNG.Float64()*0.24

	s.Targets = append(s.Targets, IntroTarget{
		ID:    s.NextTargetID,
		Type:  chosen.targetType,
		Text:  chosen.text,
		X:     x,
		Y:     y,
		Speed: speed,
		Width: textWidth(chosen.text),
		Color: chosen.color,
		MaxHP: 1,
		HP:    1,
	})
}

func (s *IntroState) spawnHazardAsteroid() {
	s.NextTargetID++
	playTop := 4
	playBottom := s.Height - 4
	yRange := playBottom - playTop
	if yRange < 1 {
		yRange = 1
	}
	y := float64(playTop + s.RNG.Intn(yRange))

	type hazardDef struct {
		name  string
		hp    int
		color lipgloss.Color
	}

	catalog := []hazardDef{
		{"🚨 FRAUD", 10, lipgloss.Color("#EF4444")},
		{"⚠ INSIDER TRADING", 12, lipgloss.Color("#F59E0B")},
		{"💣 MATERIAL WEAKNESS", 8, lipgloss.Color("#F43F5E")},
		{"💸 PONZI SCHEME", 14, lipgloss.Color("#A855F7")},
	}

	chosen := catalog[s.RNG.Intn(len(catalog))]
	text := fmt.Sprintf("[%s: HP %d/%d]", chosen.name, chosen.hp, chosen.hp)
	speed := 0.20 + s.RNG.Float64()*0.12

	s.Targets = append(s.Targets, IntroTarget{
		ID:         s.NextTargetID,
		Type:       TargetAsteroid,
		Text:       text,
		X:          float64(s.Width - 2),
		Y:          y,
		Speed:      speed,
		Width:      textWidth(text),
		Color:      chosen.color,
		MaxHP:      chosen.hp,
		HP:         chosen.hp,
		IsHazard:   true,
		HazardName: chosen.name,
	})
}

func (s *IntroState) spawnExplosion(x, y float64) {
	chars := []rune{'💥', '*', '+', 'x', 'o', '·', '°'}
	colors := []lipgloss.Color{
		lipgloss.Color("#FFFFFF"),
		lipgloss.Color("#FFE600"),
		lipgloss.Color("#FF5500"),
		lipgloss.Color("#EF4444"),
		lipgloss.Color("#94A3B8"),
	}

	count := 10 + s.RNG.Intn(6)
	for i := 0; i < count; i++ {
		angle := s.RNG.Float64() * 2 * math.Pi
		speed := 0.4 + s.RNG.Float64()*1.2
		vx := math.Cos(angle) * speed * 1.6
		vy := math.Sin(angle) * speed * 0.8
		ch := chars[s.RNG.Intn(len(chars))]
		col := colors[s.RNG.Intn(len(colors))]

		s.Particles = append(s.Particles, IntroParticle{
			X:     x,
			Y:     y,
			VX:    vx,
			VY:    vy,
			Life:  6 + s.RNG.Intn(6),
			Char:  ch,
			Color: col,
		})
	}
}

func (s *IntroState) spawnSpark(x, y float64) {
	chars := []rune{'*', '+', '·', '°'}
	for i := 0; i < 4; i++ {
		angle := s.RNG.Float64() * 2 * math.Pi
		spd := 0.3 + s.RNG.Float64()*0.6
		s.Particles = append(s.Particles, IntroParticle{
			X:     x,
			Y:     y,
			VX:    math.Cos(angle) * spd * 1.2,
			VY:    math.Sin(angle) * spd * 0.6,
			Life:  4 + s.RNG.Intn(3),
			Char:  chars[s.RNG.Intn(len(chars))],
			Color: lipgloss.Color("#FFE600"),
		})
	}
}

func (s *IntroState) spawnCallout(x, y float64, custom ...any) {
	if len(custom) >= 2 {
		txt, okTxt := custom[0].(string)
		col, okCol := custom[1].(lipgloss.Color)
		if okTxt && okCol {
			s.Callouts = append(s.Callouts, IntroCallout{
				Text:  txt,
				X:     x,
				Y:     y - 0.5,
				Life:  16,
				Color: col,
			})
			return
		}
	}

	callouts := []string{
		"+1,000 Dr!",
		"BALANCED!",
		"AUDITED!",
		"RECONCILED!",
		"DEBITS = CREDITS!",
		"GAAP VERIFIED!",
		"NET INCOME ↑",
		"T-ACCOUNT CLOSED!",
		"ZERO VARIANCE!",
		"CLEAN OPINION!",
		"EQUATION BALANCED!",
	}

	colors := []lipgloss.Color{
		lipgloss.Color("#22C55E"),
		lipgloss.Color("#38BDF8"),
		lipgloss.Color("#FACC15"),
		lipgloss.Color("#A855F7"),
	}

	text := callouts[s.RNG.Intn(len(callouts))]
	col := colors[s.RNG.Intn(len(colors))]

	s.Callouts = append(s.Callouts, IntroCallout{
		Text:  text,
		X:     x,
		Y:     y - 0.5,
		Life:  14,
		Color: col,
	})
}

// introCell holds single-character terminal rendering data.
type introCell struct {
	Char  rune
	Color lipgloss.Color
	Bold  bool
	Dim   bool
}

// Render compiles the entire animation frame into a formatted ANSI string.
func (s *IntroState) Render() string {
	w := s.Width
	h := s.Height
	if w < 50 {
		w = 50
	}
	if h < 16 {
		h = 16
	}

	// 1. Allocate 2D Screen Grid
	grid := make([][]introCell, h)
	for y := 0; y < h; y++ {
		grid[y] = make([]introCell, w)
		for x := 0; x < w; x++ {
			grid[y][x] = introCell{Char: ' '}
		}
	}

	// 2. Draw Parallax Starfield
	starColors := []lipgloss.Color{
		lipgloss.Color("#334155"),
		lipgloss.Color("#0284C7"),
		lipgloss.Color("#F8FAFC"),
	}
	starRunes := [][]rune{
		{'.', '·'},
		{'+', '*'},
		{'-', '~', '►'},
	}

	for _, star := range s.Stars {
		sx := int(star.X)
		sy := int(star.Y)
		if sx >= 0 && sx < w && sy >= 3 && sy < h-3 {
			runes := starRunes[star.Layer]
			r := runes[s.TickCount%len(runes)]
			grid[sy][sx] = introCell{
				Char:  r,
				Color: starColors[star.Layer],
				Bold:  star.Layer == 2,
				Dim:   star.Layer == 0,
			}
		}
	}

	// 3. Draw Shockwave Rings
	for _, sw := range s.Shockwaves {
		cx := int(sw.X)
		cy := int(sw.Y)
		r := sw.Radius
		for ang := 0.0; ang < 2*math.Pi; ang += 0.2 {
			px := cx + int(math.Cos(ang)*r*1.8)
			py := cy + int(math.Sin(ang)*r*0.9)
			if px >= 0 && px < w && py >= 3 && py < h-3 {
				grid[py][px] = introCell{
					Char:  '○',
					Color: lipgloss.Color("#00FFFF"),
					Bold:  true,
				}
			}
		}
	}

	// 4. Draw Incoming Accounting Targets & Hazards
	for _, t := range s.Targets {
		if ty := int(t.Y); ty >= 3 && ty < h-3 {
			putText(grid, int(t.X), ty, t.Text, t.Color, 0, w)
		}
	}

	// 5. Draw Lasers (Blaster Bolts)
	laserColor := lipgloss.Color("#00FFFF")
	if (s.TickCount/2)%2 == 0 {
		laserColor = lipgloss.Color("#FFE600")
	}
	for _, l := range s.Lasers {
		lx := int(l.X)
		ly := int(l.Y)
		if ly >= 3 && ly < h-3 {
			boltRunes := []rune("══►")
			for i, r := range boltRunes {
				col := lx + i
				if col >= 0 && col < w {
					grid[ly][col] = introCell{
						Char:  r,
						Color: laserColor,
						Bold:  true,
					}
				}
			}
		}
	}

	// 6. Draw Debris & Explosion Particles
	for _, p := range s.Particles {
		if py := int(p.Y); py >= 3 && py < h-3 {
			putText(grid, int(p.X), py, string(p.Char), p.Color, 0, w)
		}
	}

	// 7. Draw Floating Callouts
	for _, c := range s.Callouts {
		if cy := int(c.Y); cy >= 3 && cy < h-3 {
			putText(grid, int(c.X), cy, c.Text, c.Color, 0, w)
		}
	}

	// 8. Draw Spaceship (unless destroyed in game over)
	if !s.GameOver {
		s.drawShip(grid, w, h)
	}

	// 9. Draw HUD Header & Footer Frames
	s.drawHUD(grid, w, h)
	s.drawFooter(grid, w, h)

	// 10. Draw Modals (Hall of Fame, Initials Entry, Year End, Game Over, Pause)
	switch {
	case s.ShowLeaderboard:
		s.drawLeaderboardModal(grid, w, h)
	case s.InitialsEntryActive:
		s.drawInitialsModal(grid, w, h)
	case s.YearComplete:
		s.drawYearEndModal(grid, w, h)
	case s.GameOver:
		s.drawGameOverModal(grid, w, h)
	case s.Paused:
		s.drawPauseModal(grid, w, h)
	}

	// 11. Convert Grid to Final ANSI Output
	return s.renderGrid(grid, w, h)
}

// arwingFuselage is shared by every roll frame; only the wings rotate.
const arwingFuselage = "TTT[≡]=<▓▓(●)▓▓▓▓▓▓▓▓▀▀▀▀═►"

// arwingFrames holds the 5-row Arwing silhouette for each roll angle (rows are
// ship-relative -2..+2). Banking swings one wing into full planform while the
// opposite wing rotates toward edge-on, with the G-diffuser blades following.
var arwingFrames = map[ShipRoll][5]string{
	RollHardUp: {
		"            ▁▁▁▁▁▁▁▁╱▌►",
		"      ▁▂▄▆▓▓▓▓▓▓▓▓▓╱▌",
		arwingFuselage,
		"        ▔▔▀▀▀▀▀▀▀▀▀═►",
		"",
	},
	RollSoftUp: {
		"                ▁▁▁╱►",
		"      ▁▂▄▆▓▓▓▓▓▓▓▓╱",
		arwingFuselage,
		"        ▔▀▀░░░░░░╲",
		"                 ╲►",
	},
	RollLevel: {
		"",
		"       ▁▂▄▄▓▓▓▓▓▓▓╱►",
		arwingFuselage,
		"       ▔▀▀▀▓▓▓▓▓▓▓╲►",
		"",
	},
	RollSoftDown: {
		"                 ╱►",
		"        ▁▄▄░░░░░░╱",
		arwingFuselage,
		"      ▔▀▀▀▓▓▓▓▓▓▓▓╲",
		"                ▔▔▔╲►",
	},
	RollHardDown: {
		"",
		"        ▁▁▄▄▄▄▄▄▄▄▄═►",
		arwingFuselage,
		"      ▔▀▀▀▓▓▓▓▓▓▓▓▓╲▌",
		"            ▔▔▔▔▔▔▔▔╲▌►",
	},
}

func arwingRuneColor(r rune) lipgloss.Color {
	switch r {
	case '▓':
		return lipgloss.Color("#E2E8F0") // silver hull plating
	case '░':
		return lipgloss.Color("#64748B") // wing underside in shadow
	case '▁', '▂', '▄', '▆', '▀', '▔':
		return lipgloss.Color("#94A3B8") // hull contour shading
	case '╱', '╲', '▌':
		return lipgloss.Color("#00FFFF") // G-diffuser blades
	case '►', '═':
		return lipgloss.Color("#38BDF8") // laser cannons / nose
	case '●':
		return lipgloss.Color("#FFD700") // pilot canopy
	case '(', ')':
		return lipgloss.Color("#38BDF8")
	default:
		return lipgloss.Color("#475569") // engine housing
	}
}

func (s *IntroState) drawShip(grid [][]introCell, w, h int) {
	shipCenterY := int(s.ShipY)
	shipX := int(s.ShipX)

	flameColors := []lipgloss.Color{
		lipgloss.Color("#FF3B30"),
		lipgloss.Color("#FF9500"),
		lipgloss.Color("#FFCC00"),
		lipgloss.Color("#00E5FF"),
	}
	fColor := flameColors[(s.TickCount/2)%len(flameColors)]

	thrusterRunes := [][]rune{
		[]rune("»»»"),
		[]rune("░▒▓"),
		[]rune("≡≡>"),
		[]rune("*=>"),
	}
	tFlame := thrusterRunes[(s.TickCount/3)%len(thrusterRunes)]

	frame := arwingFrames[s.ShipRoll]
	for rIdx, text := range frame {
		targetY := shipCenterY + rIdx - 2
		if targetY < 3 || targetY >= h-3 {
			continue
		}
		for i, ch := range []rune(text) {
			col := shipX + i
			if ch == ' ' || col < 0 || col >= w {
				continue
			}
			colr := arwingRuneColor(ch)
			if ch == 'T' {
				ch, colr = tFlame[i], fColor
			}
			grid[targetY][col] = introCell{Char: ch, Color: colr, Bold: true}
		}
	}
}

func renderBar(val int) string {
	filled := val / 10
	if filled < 0 {
		filled = 0
	}
	if filled > 10 {
		filled = 10
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", 10-filled)
}

// textWidth returns the terminal display width of s.
func textWidth(s string) int {
	return runewidth.StringWidth(s)
}

// putCells writes text starting at column x of row y, advancing by each rune's
// display width so emoji keep the row aligned. The second half of a wide rune
// becomes a continuation cell (Char 0). Only columns in [lo, hi) are written.
func putCells(grid [][]introCell, x, y int, text string, lo, hi int, cell func(i int, r rune) introCell) {
	if y < 0 || y >= len(grid) {
		return
	}
	row := grid[y]
	lo = max(lo, 0)
	hi = min(hi, len(row))
	col := x
	for i, r := range []rune(text) {
		rw := runewidth.RuneWidth(r)
		if rw == 0 {
			continue
		}
		if col >= lo && col+rw <= hi {
			c := cell(i, r)
			c.Char = r
			row[col] = c
			if rw == 2 {
				row[col+1] = introCell{Color: c.Color}
			}
		}
		col += rw
	}
}

// putText writes bold single-color text; see putCells.
func putText(grid [][]introCell, x, y int, text string, color lipgloss.Color, lo, hi int) {
	putCells(grid, x, y, text, lo, hi, func(int, rune) introCell {
		return introCell{Color: color, Bold: true}
	})
}

// centerX returns the column that centers text within width w (minimum 2).
func centerX(text string, w int) int {
	return max(2, (w-textWidth(text))/2)
}

func (s *IntroState) drawHUD(grid [][]introCell, w, h int) {
	cyan := lipgloss.Color("#38BDF8")
	emerald := lipgloss.Color("#10B981")
	gold := lipgloss.Color("#F59E0B")
	white := lipgloss.Color("#F8FAFC")
	red := lipgloss.Color("#EF4444")
	borderCol := lipgloss.Color("#334155")

	// Line 0: Title Banner + All-Time High Score
	title := " 🚀 ACCOUNTUTOR 9000 // ORBITAL AUDIT DEFENSE "
	titleColor := cyan
	if s.Score > s.HighScore && s.Score > 0 {
		title = fmt.Sprintf(" 🏆 NEW ALL-TIME HIGH SCORE: %d! 🏆 ", s.Score)
		titleColor = gold
	}
	top := title + fmt.Sprintf(" │ ALL-TIME HIGH: [%s] %d ", s.HighScoreHolder, s.HighScore)
	putText(grid, centerX(top, w), 0, top, titleColor, 0, w)

	// Line 1: Flight, Campaign & Ordnance Stats
	pilotStatus := "DEMO"
	if s.ManualMode {
		pilotStatus = "MANUAL"
		if s.AutoFire {
			pilotStatus = "MANUAL+AUTOFIRE"
		}
	}
	bombsStr := strings.Repeat("💣", max(0, s.BombsRemaining))
	if bombsStr == "" {
		bombsStr = "EMPTY"
	}
	stats := fmt.Sprintf(" %s │ FY%d Q%d %2ds │ SCORE: %05d x%d │ SPEED: %.1fx │ BOMBS: %s ",
		pilotStatus, s.Year, s.Quarter(), s.QuarterSecondsLeft(), s.Score, s.Year, s.SpeedFactor, bombsStr)
	putCells(grid, centerX(stats, w), 1, stats, 0, w, func(_ int, r rune) introCell {
		c := white
		if r >= '0' && r <= '9' {
			c = gold
		}
		return introCell{Color: c, Bold: true}
	})

	// Line 2: Dual Audit Health Meters (Shields & Global)
	meterColor := func(hp int) lipgloss.Color {
		switch {
		case hp <= 30:
			return red
		case hp <= 60:
			return gold
		}
		return emerald
	}
	shieldColor, globalColor := meterColor(s.ShieldHP), meterColor(s.GlobalHP)
	shieldPart := fmt.Sprintf(" INT AUDIT [SHIELD]: [%s] %3d%%  │", renderBar(s.ShieldHP), s.ShieldHP)
	health := shieldPart + fmt.Sprintf("  EXT AUDIT [GLOBAL]: [%s] %3d%% ", renderBar(s.GlobalHP), s.GlobalHP)
	split := len([]rune(shieldPart))
	putCells(grid, centerX(health, w), 2, health, 0, w, func(i int, r rune) introCell {
		c := white
		if r == '█' || (r >= '0' && r <= '9') {
			c = globalColor
			if i < split {
				c = shieldColor
			}
		}
		return introCell{Color: c, Bold: true}
	})

	// Line 3: Divider
	for col := 0; col < w; col++ {
		grid[3][col] = introCell{Char: '─', Color: borderCol}
	}
}

func (s *IntroState) drawFooter(grid [][]introCell, w, h int) {
	borderCol := lipgloss.Color("#334155")
	gold := lipgloss.Color("#F59E0B")
	white := lipgloss.Color("#F8FAFC")
	cyan := lipgloss.Color("#38BDF8")
	slate := lipgloss.Color("#94A3B8")

	for col := 0; col < w; col++ {
		grid[h-3][col] = introCell{Char: '─', Color: borderCol}
	}

	promptColors := []lipgloss.Color{gold, cyan, white}
	pColor := promptColors[(s.TickCount/10)%len(promptColors)]
	var prompt, inst string
	if s.ManualMode {
		autoFire := "OFF"
		if s.AutoFire {
			autoFire = "ON"
		}
		prompt = fmt.Sprintf(">>>  HOLD [F]: MACHINE GUN  •  [G]: AUTO-FIRE %s  •  [B]: BOMB (%d)  <<<", autoFire, s.BombsRemaining)
		inst = "[W/S/↑/↓] Steer • [P] Pause • [L] Scores • [Enter] Drills • Grab ✚ to repair"
	} else {
		prompt = ">>>  PRESS ENTER TO START ACCOUNTING DRILLS  <<<   [ESC: SKIP]"
		inst = "[W/S/↑/↓] Fly • [F] Fire • [P] Pause • [L] Scores • Survive 4 quarters to win"
	}
	putText(grid, centerX(prompt, w), h-2, prompt, pColor, 0, w)
	putCells(grid, centerX(inst, w), h-1, inst, 0, w, func(_ int, r rune) introCell {
		switch r {
		case '[', ']':
			return introCell{Color: cyan}
		case '•':
			return introCell{Color: gold}
		}
		return introCell{Color: slate}
	})
}

type modalLine struct {
	text  string
	color lipgloss.Color
}

// drawModal draws a double-line box centered in the play field with each line
// centered inside it.
func drawModal(grid [][]introCell, w, h, boxW int, border lipgloss.Color, lines []modalLine) {
	boxW = min(boxW, w-4)
	boxH := len(lines) + 2
	startX := (w - boxW) / 2
	startY := max(3, (h-boxH)/2)

	for dy := 0; dy < boxH; dy++ {
		rowY := startY + dy
		if rowY >= h-3 {
			break
		}
		for dx := 0; dx < boxW; dx++ {
			colX := startX + dx
			if colX < 0 || colX >= w {
				continue
			}
			ch := ' '
			switch {
			case dy == 0 && dx == 0:
				ch = '╔'
			case dy == 0 && dx == boxW-1:
				ch = '╗'
			case dy == boxH-1 && dx == 0:
				ch = '╚'
			case dy == boxH-1 && dx == boxW-1:
				ch = '╝'
			case dy == 0 || dy == boxH-1:
				ch = '═'
			case dx == 0 || dx == boxW-1:
				ch = '║'
			}
			grid[rowY][colX] = introCell{Char: ch, Color: border, Bold: true}
		}

		if li := dy - 1; li >= 0 && li < len(lines) && lines[li].text != "" {
			text := lines[li].text
			putText(grid, startX+(boxW-textWidth(text))/2, rowY, text, lines[li].color, startX+1, startX+boxW-1)
		}
	}
}

func (s *IntroState) drawGameOverModal(grid [][]introCell, w, h int) {
	red := lipgloss.Color("#EF4444")
	gold := lipgloss.Color("#F59E0B")
	cyan := lipgloss.Color("#38BDF8")
	white := lipgloss.Color("#FFFFFF")
	emerald := lipgloss.Color("#10B981")

	border, title := red, "🚨 AUDIT FAILURE // MISSION TERMINATED 🚨"
	if s.Victory {
		border, title = emerald, "🏆 MISSION COMPLETE // BOOKS CLOSED 🏆"
	}
	drawModal(grid, w, h, 76, border, []modalLine{
		{title, border},
		{"", white},
		{s.GameOverReason, gold},
		{fmt.Sprintf("Final Score: %d  │  Entities Audited: %d  │  Reached FY%d Q%d", s.Score, s.BlastedCount, s.Year, s.Quarter()), white},
		{"", white},
		{"[R]: Retry Flight   •   [L]: High Scores   •   [Enter / Esc]: Begin Drills", cyan},
	})
}

func (s *IntroState) drawYearEndModal(grid [][]introCell, w, h int) {
	gold := lipgloss.Color("#F59E0B")
	cyan := lipgloss.Color("#38BDF8")
	white := lipgloss.Color("#FFFFFF")
	emerald := lipgloss.Color("#10B981")

	drawModal(grid, w, h, 76, gold, []modalLine{
		{fmt.Sprintf("🏆 FISCAL YEAR %d CLOSED // UNQUALIFIED OPINION 🏆", s.Year), gold},
		{"", white},
		{fmt.Sprintf("Year-end bonus: +%d  (remaining shield & integrity x FY%d)", s.YearBonus, s.Year), emerald},
		{fmt.Sprintf("Score: %d  │  Entities Audited: %d", s.Score, s.BlastedCount), white},
		{"", white},
		{fmt.Sprintf("CONTINUE INTO FISCAL YEAR %d? Faster flight, x%d points, partial repair.", s.Year+1, s.Year+1), white},
		{"[Y]: Continue   •   [N]: Retire & bank your score", cyan},
	})
}

func (s *IntroState) drawPauseModal(grid [][]introCell, w, h int) {
	cyan := lipgloss.Color("#38BDF8")
	white := lipgloss.Color("#FFFFFF")
	slate := lipgloss.Color("#94A3B8")

	drawModal(grid, w, h, 60, cyan, []modalLine{
		{"❚❚  FLIGHT PAUSED  ❚❚", cyan},
		{"", white},
		{fmt.Sprintf("FY%d · Q%d · %ds left in quarter · Score %d", s.Year, s.Quarter(), s.QuarterSecondsLeft(), s.Score), slate},
		{"", white},
		{"[P]: Resume   •   [L]: High Scores", white},
		{"[Enter / Esc]: Accounting Drills   •   [Q]: Quit", white},
	})
}

func (s *IntroState) drawInitialsModal(grid [][]introCell, w, h int) {
	gold := lipgloss.Color("#F59E0B")
	emerald := lipgloss.Color("#10B981")
	cyan := lipgloss.Color("#38BDF8")
	white := lipgloss.Color("#FFFFFF")

	lines := []modalLine{
		{fmt.Sprintf("🏆 NEW ALL-TIME HIGH SCORE: %d! 🏆", s.Score), gold},
		{"ENTER PILOT INITIALS FOR THE AUDIT HALL OF FAME:", white},
		{"", white},
		{"", white}, // initials slots
		{"", white}, // cursor
		{"[↑/↓] or [A-Z]: Select Letter  •  [Enter/Space]: Next/Save", cyan},
	}
	drawModal(grid, w, h, 66, gold, lines)

	// 3 Slots for Initials: [ A ] [ B ] [ C ], on the rows drawModal left blank.
	boxW := min(66, w-4)
	startX := (w - boxW) / 2
	startY := max(3, (h-(len(lines)+2))/2)
	slotsY := startY + 4
	if slotsY+1 >= h-3 {
		return
	}
	slotStartX := startX + (boxW-19)/2
	for slot := 0; slot < 3; slot++ {
		sColor := cyan
		if slot == s.InitialsCursor {
			sColor = emerald
			grid[slotsY+1][slotStartX+slot*7+2] = introCell{Char: '▲', Color: emerald, Bold: true}
		}
		putText(grid, slotStartX+slot*7, slotsY, fmt.Sprintf("[ %c ]", s.Initials[slot]), sColor, 0, w)
	}
}

func (s *IntroState) drawLeaderboardModal(grid [][]introCell, w, h int) {
	gold := lipgloss.Color("#F59E0B")
	cyan := lipgloss.Color("#38BDF8")
	white := lipgloss.Color("#FFFFFF")
	slate := lipgloss.Color("#94A3B8")
	emerald := lipgloss.Color("#10B981")

	lines := []modalLine{
		{"★  AUDIT HALL OF FAME  ★", gold},
		{"", white},
		{fmt.Sprintf("%-4s  %-5s  %9s  %7s  %6s", "RANK", "PILOT", "SCORE", "AUDITED", "TIME"), slate},
	}
	if len(s.Leaderboard) == 0 {
		lines = append(lines, modalLine{"No flights logged yet. Be the first!", white})
	}
	for i, e := range s.Leaderboard {
		c := white
		if i == 0 {
			c = gold
		}
		if s.ScoreSaved && e.Score == s.Score && e.Initials == s.GetInitialsString() {
			c = emerald
		}
		lines = append(lines, modalLine{fmt.Sprintf("%-4s  %-5s  %9d  %7d  %5ds", fmt.Sprintf("#%d", i+1), e.Initials, e.Score, e.BlastedCount, e.SurvivalSeconds), c})
	}
	lines = append(lines, modalLine{"", white}, modalLine{"[L / Esc]: Close  •  [R]: New Flight  •  [Enter]: Drills", cyan})
	drawModal(grid, w, h, 64, gold, lines)
}

func (s *IntroState) renderGrid(grid [][]introCell, w, h int) string {
	var sb strings.Builder
	for y := 0; y < h; y++ {
		row := grid[y]
		var runChars []rune
		var curColor lipgloss.Color
		var curBold bool
		var curDim bool
		hasRun := false

		flushRun := func() {
			if len(runChars) == 0 {
				return
			}
			txt := string(runChars)
			if curColor == "" && !curBold && !curDim {
				sb.WriteString(txt)
			} else {
				st := lipgloss.NewStyle()
				if curColor != "" {
					st = st.Foreground(curColor)
				}
				if curBold {
					st = st.Bold(true)
				}
				if curDim {
					st = st.Faint(true)
				}
				sb.WriteString(st.Render(txt))
			}
			runChars = runChars[:0]
		}

		for x := 0; x < w; x++ {
			c := row[x]
			// Keep every row exactly w columns: a wide rune consumes its
			// continuation cell, and one whose half was overdrawn (or an
			// orphaned continuation) renders as a blank.
			switch {
			case c.Char == 0:
				c.Char = ' '
			case runewidth.RuneWidth(c.Char) == 2:
				if x+1 < w && row[x+1].Char == 0 {
					x++
				} else {
					c.Char = ' '
				}
			}
			if !hasRun {
				curColor = c.Color
				curBold = c.Bold
				curDim = c.Dim
				runChars = append(runChars, c.Char)
				hasRun = true
			} else if c.Color == curColor && c.Bold == curBold && c.Dim == curDim {
				runChars = append(runChars, c.Char)
			} else {
				flushRun()
				curColor = c.Color
				curBold = c.Bold
				curDim = c.Dim
				runChars = append(runChars, c.Char)
			}
		}
		flushRun()
		if y < h-1 {
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}
