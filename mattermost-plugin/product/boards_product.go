// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package product

import (
	"errors"
	"fmt"

	"github.com/mattermost/focalboard/mattermost-plugin/server/boards"
	"github.com/mattermost/focalboard/server/model"

	"github.com/mattermost/mattermost-server/v6/app"
	"github.com/mattermost/mattermost-server/v6/einterfaces"
	mm_model "github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
	"github.com/mattermost/mattermost-server/v6/product"
	"github.com/mattermost/mattermost-server/v6/shared/mlog"
)

const (
	boardsProductName = "boards"
	boardsProductID   = "com.mattermost.boards"
)

var errServiceTypeAssert = errors.New("type assertion failed")

func init() {
	app.RegisterProduct(boardsProductName, app.ProductManifest{
		Initializer: newBoardsProduct,
		Dependencies: map[app.ServiceKey]struct{}{
			app.TeamKey:          {},
			app.ChannelKey:       {},
			app.UserKey:          {},
			app.PostKey:          {},
			app.BotKey:           {},
			app.ClusterKey:       {},
			app.ConfigKey:        {},
			app.LogKey:           {},
			app.LicenseKey:       {},
			app.FilestoreKey:     {},
			app.FileInfoStoreKey: {},
			app.RouterKey:        {},
			app.CloudKey:         {},
			app.KVStoreKey:       {},
			app.StoreKey:         {},
			app.SystemKey:        {},
			app.PreferencesKey:   {},
		},
	})
}

type boardsProduct struct {
	teamService          product.TeamService
	channelService       product.ChannelService
	userService          product.UserService
	postService          product.PostService
	permissionsService   product.PermissionService
	botService           product.BotService
	clusterService       product.ClusterService
	configService        product.ConfigService
	logger               mlog.LoggerIFace
	licenseService       product.LicenseService
	filestoreService     product.FilestoreService
	fileInfoStoreService product.FileInfoStoreService
	routerService        product.RouterService
	cloudService         product.CloudService
	kvStoreService       product.KVStoreService
	storeService         product.StoreService
	systemService        product.SystemService
	preferencesService   product.PreferencesService
	hooksService         product.HooksService

	boardsApp *boards.BoardsApp
}

func newBoardsProduct(_ *app.Server, services map[app.ServiceKey]interface{}) (app.Product, error) {
	boardsProduct := &boardsProduct{}

	for key, service := range services {
		switch key {
		case app.TeamKey:
			teamService, ok := service.(product.TeamService)
			if !ok {
				return nil, fmt.Errorf("invalid service key '%s': %w", key, errServiceTypeAssert)
			}
			boardsProduct.teamService = teamService
		case app.ChannelKey:
			channelService, ok := service.(product.ChannelService)
			if !ok {
				return nil, fmt.Errorf("invalid service key '%s': %w", key, errServiceTypeAssert)
			}
			boardsProduct.channelService = channelService
		case app.UserKey:
			userService, ok := service.(product.UserService)
			if !ok {
				return nil, fmt.Errorf("invalid service key '%s': %w", key, errServiceTypeAssert)
			}
			boardsProduct.userService = userService
		case app.PostKey:
			postService, ok := service.(product.PostService)
			if !ok {
				return nil, fmt.Errorf("invalid service key '%s': %w", key, errServiceTypeAssert)
			}
			boardsProduct.postService = postService
		case app.PermissionsKey:
			permissionsService, ok := service.(product.PermissionService)
			if !ok {
				return nil, fmt.Errorf("invalid service key '%s': %w", key, errServiceTypeAssert)
			}
			boardsProduct.permissionsService = permissionsService
		case app.BotKey:
			botService, ok := service.(product.BotService)
			if !ok {
				return nil, fmt.Errorf("invalid service key '%s': %w", key, errServiceTypeAssert)
			}
			boardsProduct.botService = botService
		case app.ClusterKey:
			clusterService, ok := service.(product.ClusterService)
			if !ok {
				return nil, fmt.Errorf("invalid service key '%s': %w", key, errServiceTypeAssert)
			}
			boardsProduct.clusterService = clusterService
		case app.ConfigKey:
			configService, ok := service.(product.ConfigService)
			if !ok {
				return nil, fmt.Errorf("invalid service key '%s': %w", key, errServiceTypeAssert)
			}
			boardsProduct.configService = configService
		case app.LogKey:
			logger, ok := service.(mlog.LoggerIFace)
			if !ok {
				return nil, fmt.Errorf("invalid service key '%s': %w", key, errServiceTypeAssert)
			}
			boardsProduct.logger = logger.With(mlog.String("product", boardsProductName))
		case app.LicenseKey:
			licenseService, ok := service.(product.LicenseService)
			if !ok {
				return nil, fmt.Errorf("invalid service key '%s': %w", key, errServiceTypeAssert)
			}
			boardsProduct.licenseService = licenseService
		case app.FilestoreKey:
			filestoreService, ok := service.(product.FilestoreService)
			if !ok {
				return nil, fmt.Errorf("invalid service key '%s': %w", key, errServiceTypeAssert)
			}
			boardsProduct.filestoreService = filestoreService
		case app.FileInfoStoreKey:
			fileInfoStoreService, ok := service.(product.FileInfoStoreService)
			if !ok {
				return nil, fmt.Errorf("invalid service key '%s': %w", key, errServiceTypeAssert)
			}
			boardsProduct.fileInfoStoreService = fileInfoStoreService
		case app.RouterKey:
			routerService, ok := service.(product.RouterService)
			if !ok {
				return nil, fmt.Errorf("invalid service key '%s': %w", key, errServiceTypeAssert)
			}
			boardsProduct.routerService = routerService
		case app.CloudKey:
			cloudService, ok := service.(product.CloudService)
			if !ok {
				return nil, fmt.Errorf("invalid service key '%s': %w", key, errServiceTypeAssert)
			}
			boardsProduct.cloudService = cloudService
		case app.KVStoreKey:
			kvStoreService, ok := service.(product.KVStoreService)
			if !ok {
				return nil, fmt.Errorf("invalid service key '%s': %w", key, errServiceTypeAssert)
			}
			boardsProduct.kvStoreService = kvStoreService
		case app.StoreKey:
			storeService, ok := service.(product.StoreService)
			if !ok {
				return nil, fmt.Errorf("invalid service key '%s': %w", key, errServiceTypeAssert)
			}
			boardsProduct.storeService = storeService
		case app.SystemKey:
			systemService, ok := service.(product.SystemService)
			if !ok {
				return nil, fmt.Errorf("invalid service key '%s': %w", key, errServiceTypeAssert)
			}
			boardsProduct.systemService = systemService
		case app.PreferencesKey:
			preferencesService, ok := service.(product.PreferencesService)
			if !ok {
				return nil, fmt.Errorf("invalid service key '%s': %w", key, errServiceTypeAssert)
			}
			boardsProduct.preferencesService = preferencesService
		case app.HooksKey:
			hooksService, ok := service.(product.HooksService)
			if !ok {
				return nil, fmt.Errorf("invalid service key '%s': %w", key, errServiceTypeAssert)
			}
			boardsProduct.hooksService = hooksService
		}
	}

	if !boardsProduct.configService.Config().FeatureFlags.BoardsProduct {
		boardsProduct.logger.Info("Boards product disabled via feature flag")
		return boardsProduct, nil
	}

	adapter := newServiceAPIAdapter(boardsProduct)
	boardsApp, err := boards.NewBoardsApp(adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create Boards service: %w", err)
	}

	model.LogServerInfo(boardsProduct.logger)

	boardsProduct.boardsApp = boardsApp

	return boardsProduct, nil
}

func (bp *boardsProduct) Start() error {
	if !bp.configService.Config().FeatureFlags.BoardsProduct {
		bp.logger.Info("Boards product disabled via feature flag")
		return nil
	}

	bp.logger.Info("Starting boards service")

	if err := bp.hooksService.RegisterHooks(boardsProductName, bp); err != nil {
		return fmt.Errorf("failed to register hooks: %w", err)
	}

	if err := bp.boardsApp.Start(); err != nil {
		return fmt.Errorf("failed to start Boards service: %w", err)
	}

	return nil
}

func (bp *boardsProduct) Stop() error {
	bp.logger.Info("Stopping boards service")

	if bp.boardsApp == nil {
		return nil
	}

	if err := bp.boardsApp.Stop(); err != nil {
		return fmt.Errorf("error while stopping Boards service: %w", err)
	}

	return nil
}

//
// These callbacks are called by the suite automatically
//

func (bp *boardsProduct) OnConfigurationChange() error {
	if bp.boardsApp == nil {
		return nil
	}
	return bp.boardsApp.OnConfigurationChange()
}

func (bp *boardsProduct) OnWebSocketConnect(webConnID, userID string) {
	if bp.boardsApp == nil {
		return
	}
	bp.boardsApp.OnWebSocketConnect(webConnID, userID)
}

func (bp *boardsProduct) OnWebSocketDisconnect(webConnID, userID string) {
	if bp.boardsApp == nil {
		return
	}
	bp.boardsApp.OnWebSocketDisconnect(webConnID, userID)
}

func (bp *boardsProduct) WebSocketMessageHasBeenPosted(webConnID, userID string, req *mm_model.WebSocketRequest) {
	if bp.boardsApp == nil {
		return
	}
	bp.boardsApp.WebSocketMessageHasBeenPosted(webConnID, userID, req)
}

func (bp *boardsProduct) OnPluginClusterEvent(ctx *plugin.Context, ev mm_model.PluginClusterEvent) {
	if bp.boardsApp == nil {
		return
	}
	bp.boardsApp.OnPluginClusterEvent(ctx, ev)
}

func (bp *boardsProduct) MessageWillBePosted(ctx *plugin.Context, post *mm_model.Post) (*mm_model.Post, string) {
	if bp.boardsApp == nil {
		return post, ""
	}
	return bp.boardsApp.MessageWillBePosted(ctx, post)
}

func (bp *boardsProduct) MessageWillBeUpdated(ctx *plugin.Context, newPost, oldPost *mm_model.Post) (*mm_model.Post, string) {
	if bp.boardsApp == nil {
		return newPost, ""
	}
	return bp.boardsApp.MessageWillBeUpdated(ctx, newPost, oldPost)
}

func (bp *boardsProduct) OnCloudLimitsUpdated(limits *mm_model.ProductLimits) {
	if bp.boardsApp == nil {
		return
	}
	bp.boardsApp.OnCloudLimitsUpdated(limits)
}

func (bp *boardsProduct) RunDataRetention(nowTime, batchSize int64) (int64, error) {
	if bp.boardsApp == nil {
		return 0, nil
	}
	return bp.boardsApp.RunDataRetention(nowTime, batchSize)
}

func (bp *boardsProduct) Exporter(cursor map[string]any, limit int) (einterfaces.ComplianceExporter, map[string]any, error) {
	if bp.boardsApp == nil {
		return nil, nil, nil
	}
	return bp.boardsApp.Exporter(cursor, limit)
}
