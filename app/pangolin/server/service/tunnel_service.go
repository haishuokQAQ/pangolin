package service

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"pangolin/app/pangolin/constant"
	"pangolin/app/pangolin/model/db"
	"pangolin/app/pangolin/utils"
	"time"
)

func (srv *Service) CreateTunnel(ctx context.Context, config *db.TunnelConfig) error {
	tunnel, err := srv.initTunnelFromConfig(ctx, config)
	if err != nil {
		return err
	}
	result := srv.addIntoTunnelMap(config, tunnel)
	if result < 0 {
		return constant.GetErrorByErrorCode(result)
	}
	go utils.SafeGoroutine(tunnel.Start, func(i interface{}) {
		srv.logger.Errorf("Error occurred when run tunner %+v", config.Id)

	})
	return nil
}

func (srv *Service) addIntoTunnelMap(config *db.TunnelConfig, sshTunnel *utils.SSHTunnel) int {
	srv.tunnelMapLock.Lock()
	defer srv.tunnelMapLock.Unlock()
	if _, ok := srv.portMap[config.LocalPort]; !ok {
		return -1
	}
	if _, ok := srv.tunnelMap[config.Id]; !ok {
		return -2
	}
	srv.portMap[config.LocalPort] = true
	srv.tunnelMap[config.Id] = sshTunnel
	srv.tunnelConfig[config.Id] = config
	return 0
}

func (srv *Service) initTunnelFromConfig(ctx context.Context, config *db.TunnelConfig) (*utils.SSHTunnel, error) {
	tunnel := &utils.SSHTunnel{
		Local: &utils.Endpoint{
			Host: config.LocalHost,
			Port: config.LocalPort,
		},
		Server: &utils.Endpoint{
			Host: config.ServerHost,
			Port: config.ServerPort,
		},
		Remote: &utils.Endpoint{
			Host: config.RemoteHost,
			Port: config.RemotePort,
		},
	}
	tunnel.Config = &ssh.ClientConfig{
		Timeout: 5 * time.Second,
		User:    config.UserName,
	}
	switch config.LogMode {
	case db.LogModePassword:
		tunnel.Config.Auth = []ssh.AuthMethod{ssh.Password(config.Password)}
		break
	case db.LogModePrivateKey:
		key, err := utils.PrivateKeyString(config.PrivateKey)
		if err != nil {
			srv.logger.WithTraceInCtx(ctx).Errorf("Fail to parse private key.Error %+v", err)
			return nil, err
		}
		tunnel.Config.Auth = []ssh.AuthMethod{key}
		break
	default:
		srv.logger.WithTraceInCtx(ctx).Errorf("Illegal log mode %+v", config.LogMode)
		return nil, errors.New(fmt.Sprintf("Illegal log mode %+v", config.LogMode))
	}
	return tunnel, nil
}

func (srv *Service) GetExistTunnel(ctx context.Context) []*db.TunnelConfig {
	srv.tunnelMapLock.RLock()
	defer srv.tunnelMapLock.RUnlock()
	result := []*db.TunnelConfig{}
	for _, config := range srv.tunnelConfig {
		result = append(result, config)
	}
	return result
}

func (srv *Service) DestroyTunnel(ctx context.Context, configId uint64) error {
	srv.tunnelMapLock.Lock()
	defer srv.tunnelMapLock.Unlock()
	if _, ok := srv.tunnelMap[configId]; !ok {
		return errors.New("Config not exist!")
	}
	tunnel := srv.tunnelMap[configId]
	delete(srv.tunnelMap, configId)
	delete(srv.tunnelConfig, configId)
	delete(srv.portMap, tunnel.Local.Port)
	tunnel.Shutdown()
	return nil
}
