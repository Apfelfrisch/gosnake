package game

type Perks map[PerkType]Perk

func (p *Perks) add(pt PerkType, usages uint16) {
	if *p == nil {
		*p = make(Perks)
	}

	(*p)[pt] = Perk{usages}
}

func (ps *Perks) Get(pt PerkType) Perk {
	if *ps == nil {
		*ps = make(Perks)
	}

	if p, ok := (*ps)[pt]; ok {
		return p
	}

	ps.add(pt, 0)

	return Perk{0}
}

func (ps *Perks) set(pt PerkType, p Perk) {
	if *ps == nil {
		*ps = make(Perks)
	}

	(*ps)[pt] = p
}

func (ps *Perks) use(pt PerkType) bool {
	p := ps.Get(pt)

	if p.Usages == 0 {
		return false
	}

	p.Usages -= 1

	ps.set(pt, p)

	return true
}

func (ps *Perks) reload(pt PerkType, usages uint16) {
	p := ps.Get(pt)

	p.Usages += usages

	ps.set(pt, p)
}

type PerkType int

func (pt PerkType) String() string {
	switch pt {
	case WalkWall:
		return "Walk Wall"
	case Dash:
		return "Stash"
	}

	return "Unkown"
}

const (
	WalkWall PerkType = 1
	Dash     PerkType = 2
)

type Perk struct {
	Usages uint16 `json:"u"`
}

func (p *Perk) reload(usages uint16) {
	p.Usages += usages
}
