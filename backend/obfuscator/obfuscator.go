package obfuscator

import (
	"context"
	"errors"
	"hyperion/backend/obfuscator/js"
	"os"
	"path/filepath"
)

type Obfuscator struct {
	ctx    context.Context
	config *Config
}

func NewObfuscator(config *Config) *Obfuscator {
	return &Obfuscator{
		config: config,
	}
}

type IObfuscation interface {
	SingleLine(code string) (string, error)
	RefactorLoopStatement(code string) (string, error)
	RefactorIfStatement(code string) (string, error)
	RenameVariable(code string) (string, error)
	RenameConstant(code string) (string, error)
	RenameFunction(code string) (string, error)
	RemoveComment(code string) (string, error)
}

func GetObfuscation(extension string) (IObfuscation, error) {
	switch extension {
	case "js":
		return js.Obfuscation{}, nil
	default:
		return nil, errors.New("unknown obfuscation extension")
	}
}

func (o *Obfuscator) SetContext(ctx context.Context) {
	o.ctx = ctx
}

func (o *Obfuscator) Obfuscate(path string) (string, error) {
	if path == "" {
		return "", errors.New("empty path")
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	content := string(b)

	extension := filepath.Ext(path)
	obfuscation, err := GetObfuscation(extension)
	if err != nil {
		return "", err
	}

	if o.config.RemoveComments {
		content, err = obfuscation.RemoveComment(content)
		if err != nil {
			return content, err
		}
	}

	if o.config.VariableName {
		content, err = obfuscation.RenameVariable(content)
		if err != nil {
			return content, err
		}
	}

	if o.config.ConstantName {
		content, err = obfuscation.RenameConstant(content)
		if err != nil {
			return content, err
		}
	}

	if o.config.FunctionName {
		content, err = obfuscation.RenameFunction(content)
		if err != nil {
			return content, err
		}
	}

	if o.config.IfStatement {
		content, err = obfuscation.RefactorIfStatement(content)
		if err != nil {
			return content, err
		}
	}

	if o.config.LoopStatement {
		content, err = obfuscation.RefactorLoopStatement(content)
		if err != nil {
			return content, err
		}
	}

	if o.config.SingleLineOutput {
		content, err = obfuscation.SingleLine(content)
		if err != nil {
			return content, err
		}
	}

	return content, nil
}
