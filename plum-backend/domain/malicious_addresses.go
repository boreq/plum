package domain

type MaliciousAddresses struct {
	hits map[string]int
}

func NewMaliciousAddresses() *MaliciousAddresses {
	return &MaliciousAddresses{
		hits: make(map[string]int),
	}
}

func (m *MaliciousAddresses) Insert(data *Data) {
	categoryData, ok := data.Categories[CategoryMalicious]
	if !ok {
		return
	}

	for remoteAddress, remoteAddressData := range categoryData.RemoteAddresses {
		m.hits[remoteAddress] += remoteAddressData.Hits
	}
}

func (m *MaliciousAddresses) IsMalicious(remoteAddress string) bool {
	return m.hits[remoteAddress] > MaliciousRequestThreshold
}
