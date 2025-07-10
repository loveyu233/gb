package captcha

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
)

// Manager 验证码管理器
type Manager struct {
	cache      CacheImpl
	generators map[CaptchaType]CaptchaInterface
	config     CaptchaConfig
	mu         sync.RWMutex
}

// NewManager 创建验证码管理器
func NewManager(cache CacheImpl, config ...CaptchaConfig) *Manager {
	cfg := DefaultConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	return &Manager{
		cache:      cache,
		generators: make(map[CaptchaType]CaptchaInterface),
		config:     cfg,
	}
}

// RegisterGenerator 注册验证码生成器
func (m *Manager) RegisterGenerator(generator CaptchaInterface) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.generators[generator.GetType()] = generator
}

// Generate 生成验证码
func (m *Manager) Generate(captchaType CaptchaType, configs ...CaptchaConfig) (*CaptchaResponse, error) {
	generator, exists := m.generators[captchaType]
	config := DefaultConfig
	if len(configs) > 0 {
		config = configs[0]
	}
	if !exists {
		return nil, fmt.Errorf("unsupported captcha type: %s", captchaType)
	}

	// 生成验证码
	response, err := generator.Generate(m.generateCacheKey(captchaType), config)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCaptchaGenerate, err)
	}

	return response, nil
}

// Verify 验证验证码
func (m *Manager) Verify(req *CaptchaVerifyRequest) error {
	captchaType, err := m.extractCaptchaType(req.CaptchaKey)
	if err != nil {
		return err
	}
	generator, exists := m.generators[captchaType]
	if !exists {
		return fmt.Errorf("unsupported captcha type: %s", captchaType)
	}

	// 验证验证码
	if err := generator.Verify(req.CaptchaKey, req.CaptchaData, req.ValidationTolerance); err != nil {
		return fmt.Errorf("%w: %v", ErrCaptchaVerify, err)
	}

	// 验证成功后删除缓存
	go func() {
		if err := m.cache.DelCaptcha(req.CaptchaKey); err != nil {
			// 记录日志而不是返回错误
			fmt.Printf("删除验证码缓存失败: %v\n", err)
		}
	}()

	return nil
}

// generateCacheKey 生成缓存键
func (m *Manager) generateCacheKey(captchaType CaptchaType) string {
	return fmt.Sprintf("%s:%s", captchaType, uuid.NewString())
}

// extractCaptchaType 从缓存键中提取验证码类型
func (m *Manager) extractCaptchaType(key string) (CaptchaType, error) {
	// 简单的键格式: type:uuid
	if len(key) < 6 {
		return "", ErrInvalidData
	}

	switch {
	case key[:5] == "click":
		return "click", nil
	case key[:6] == "rotate":
		return "rotate", nil
	case key[:5] == "slide":
		return "slide", nil
	default:
		return "", ErrInvalidData
	}
}

// GetSupportedTypes 获取支持的验证码类型
func (m *Manager) GetSupportedTypes() []CaptchaType {
	m.mu.RLock()
	defer m.mu.RUnlock()

	types := make([]CaptchaType, 0, len(m.generators))
	for t := range m.generators {
		types = append(types, t)
	}
	return types
}
