package models

import (
	"SuperBizAgent/utility/config"
	"context"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
)

func OpenAIForDeepSeekV31Think(ctx context.Context) (cm model.ToolCallingChatModel, err error) {
	cfg := &openai.ChatModelConfig{
		Model:   config.GetString("ds_think_chat_model.model"),
		APIKey:  config.GetString("ds_think_chat_model.api_key"),
		BaseURL: config.GetString("ds_think_chat_model.base_url"),
	}
	cm, err = openai.NewChatModel(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return cm, nil
}

func OpenAIForDeepSeekV3Quick(ctx context.Context) (cm model.ToolCallingChatModel, err error) {
	cfg := &openai.ChatModelConfig{
		Model:   config.GetString("ds_quick_chat_model.model"),
		APIKey:  config.GetString("ds_quick_chat_model.api_key"),
		BaseURL: config.GetString("ds_quick_chat_model.base_url"),
	}
	cm, err = openai.NewChatModel(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return cm, nil
}
