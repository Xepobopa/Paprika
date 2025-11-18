package service

import (
	"Paprika/models"
	"Paprika/process"
	"fmt"
)

func GetAllAvailableExchanges() []string {
	return []string{models.MEXC, models.KUCOIN}
}

func SpawnExchangeProcess(name string, pm *process.Manager) error {
	switch name {

	case models.MEXC:
		return pm.Spawn(name, process.SpawnMexcProcess())
	
	case models.KUCOIN:
		return pm.Spawn(name, process.SpawnKucoinProcess())
		
	default:
		// unknown exchange name
		return fmt.Errorf("provided exchange name '%s' is not supported!", name)

	}
}

func GetExchangeStats(name string, pm *process.Manager) (*process.ProcessInformation, error) {
	res, ok := pm.GetStats(name)
	if !ok {
		return nil, fmt.Errorf("failed to get stats for '%s' process, because process with this name does not exists!", name)
	}

	return res, nil
}

func StopExchangeByName(name string, pm *process.Manager) error {
	return pm.Stop(name)
}