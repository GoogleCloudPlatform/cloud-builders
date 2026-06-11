package testservices

import (
	"context"
	"fmt"

	"sigs.k8s.io/yaml"
)

// TestKubectl implements the KubectlService interface.
type TestKubectl struct {
	ApplyResponse           map[string][]error
	ApplyFromStringResponse map[string][]error
	GetResponse             map[string]map[string][]GetResponse
	GetFromStringResponse   map[string][]GetResponse
}

// StatResponse represents a response tuple for a Stat function call.
type GetResponse struct {
	Res string
	Err error
}

// Apply calls `kubectl apply -f <filename> -n <namespace>`.
func (k *TestKubectl) Apply(ctx context.Context, filename, namespace string) error {
	errors, ok := k.ApplyResponse[filename]
	if !ok {
		panic(fmt.Sprintf("ApplyResponse has no response for filename %q", filename))
	}
	if len(errors) == 0 {
		panic(fmt.Sprintf("ApplyResponse ran out of responses for filename %q", filename))
	}
	err := errors[0]
	if len(errors) == 1 {
		delete(k.ApplyResponse, filename)
	} else {
		k.ApplyResponse[filename] = k.ApplyResponse[filename][:1]
	}
	return err
}

// ApplyFromString calls `kubectl apply -f - -n <namespace> < ${configString}`.
func (k *TestKubectl) ApplyFromString(ctx context.Context, configString, namespace string) error {
	errors, ok := k.ApplyFromStringResponse[configString]
	if !ok {
		panic(fmt.Sprintf("ApplyFromStringResponse has no response for configs %q", configString))
	}
	if len(errors) == 0 {
		panic(fmt.Sprintf("ApplyFromStringResponse ran out of responses for configs %q", configString))
	}
	err := errors[0]
	if len(errors) == 1 {
		delete(k.ApplyFromStringResponse, configString)
	} else {
		k.ApplyFromStringResponse[configString] = k.ApplyFromStringResponse[configString][1:]
	}
	return err
}

// Get calls `kubectl get <kind> <name> -n <namespace> --output=<format>`.
func (k *TestKubectl) Get(ctx context.Context, kind, name, namespace, format string, ignoreNotFound bool) (string, error) {
	resp, ok := k.GetResponse[kind][name]
	if !ok {
		panic(fmt.Sprintf("GetResponse has no response for kind %q and name %q", kind, name))
	}

	if len(resp) == 0 {
		panic(fmt.Sprintf("GetResponse ran out of responses for kind %q and name %q after", kind, name))
	}
	res := resp[0].Res
	err := resp[0].Err

	if len(resp) == 1 {
		delete(k.GetResponse[kind], name)
		if len(k.GetResponse[kind]) == 0 {
			delete(k.GetResponse, kind)
		}
	} else {
		k.GetResponse[kind][name] = k.GetResponse[kind][name][1:]
	}
	return res, err
}

// GetFromString calls `kubectl get -f - -n <namespace> < ${configString}`.
func (k *TestKubectl) GetFromString(ctx context.Context, configString, namespace, format string) (string, error) {
	if k.GetFromStringResponse != nil {
		resp, ok := k.GetFromStringResponse[configString]
		if !ok {
			panic(fmt.Sprintf("GetFromStringResponse has no response for configs %q", configString))
		}
		if len(resp) == 0 {
			panic(fmt.Sprintf("GetFromStringResponse ran out of responses for configs %q", configString))
		}
		res := resp[0].Res
		err := resp[0].Err
		if len(resp) == 1 {
			delete(k.GetFromStringResponse, configString)
		} else {
			k.GetFromStringResponse[configString] = k.GetFromStringResponse[configString][1:]
		}
		return res, err
	}

	kind, name, err := objectKindAndName(configString)
	if err != nil {
		return "", err
	}
	return k.Get(ctx, kind, name, namespace, format, false)
}

func objectKindAndName(configString string) (string, string, error) {
	var obj struct {
		Kind     string `yaml:"kind"`
		Metadata struct {
			Name string `yaml:"name"`
		} `yaml:"metadata"`
	}
	if err := yaml.Unmarshal([]byte(configString), &obj); err != nil {
		return "", "", err
	}
	if obj.Kind == "" || obj.Metadata.Name == "" {
		return "", "", fmt.Errorf("config missing kind or metadata.name")
	}
	return obj.Kind, obj.Metadata.Name, nil
}
