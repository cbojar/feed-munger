package feeds

type Namespaces []*Namespace

type Namespace struct {
	key   string
	value string
}

func NewNamespace(key string, value string) *Namespace {
	return &Namespace{key: key, value: value}
}

func (ns *Namespace) Key() string {
	return ns.key
}

func (ns *Namespace) Value() string {
	return ns.value
}
