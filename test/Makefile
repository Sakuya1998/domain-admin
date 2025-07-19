# API 测试 Makefile

.PHONY: test bench coverage clean help

# 默认目标
all: test

# 运行所有测试
test:
	@echo "运行所有测试..."
	go test -v

# 运行性能基准测试
bench:
	@echo "运行性能基准测试..."
	go test -bench=. -v

# 生成测试覆盖率报告
coverage:
	@echo "生成测试覆盖率报告..."
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 运行特定测试
test-auth:
	@echo "运行认证相关测试..."
	go test -v -run "Test.*Login|TestUnauthorized|TestInvalidToken"

test-register:
	@echo "运行注册相关测试..."
	go test -v -run "Test.*Register"

test-user:
	@echo "运行用户相关测试..."
	go test -v -run "TestUser.*"

test-admin:
	@echo "运行管理员相关测试..."
	go test -v -run "Test.*Admin|TestCreate|TestUpdate|TestDelete"

test-validation:
	@echo "运行验证相关测试..."
	go test -v -run "TestValidation|TestForbidden"

# 清理生成的文件
clean:
	@echo "清理生成的文件..."
	rm -f coverage.out coverage.html

# 安装依赖
deps:
	@echo "安装测试依赖..."
	go mod tidy
	go mod download

# 显示帮助信息
help:
	@echo "可用的命令:"
	@echo "  test        - 运行所有测试"
	@echo "  bench       - 运行性能基准测试"
	@echo "  coverage    - 生成测试覆盖率报告"
	@echo "  test-auth   - 运行认证相关测试"
	@echo "  test-register - 运行注册相关测试"
	@echo "  test-user   - 运行用户相关测试"
	@echo "  test-admin  - 运行管理员相关测试"
	@echo "  test-validation - 运行验证相关测试"
	@echo "  clean       - 清理生成的文件"
	@echo "  deps        - 安装测试依赖"
	@echo "  help        - 显示此帮助信息"