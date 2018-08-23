package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"text/template"
)

var (
	configPath string
	outputPath string
)

func main() {
	flag.StringVar(&configPath, "c", "", "config file path")
	flag.StringVar(&outputPath, "o", "out.go", "output file path")
	flag.Parse()

	log.Printf("configPath: %v", configPath)
	log.Printf("outputPath: %v", outputPath)

	t, e := template.New("xxx").Parse(Template)
	if e != nil {
		log.Fatalf("new template failed, %v", e)
	}

	c, e := ReadConfigFile(configPath)
	if e != nil {
		log.Fatalf("read config file failed, %v", e)
	}

	b := new(bytes.Buffer)

	e = t.Execute(b, c)
	if e != nil {
		log.Fatalf("execute failed, %v", e)
	}

	e = ioutil.WriteFile(outputPath, b.Bytes(), 0664)
	if e != nil {
		log.Fatalf("write file failed, %v", e)
	}
}

type Config struct {
	Name    string `json:"name"`
	Plural  string `json:"plural"`
	VarName string `json:"varName"`

	IsNonNamespaced bool `json:"isNonNamespaced"`

	ImportPath string `json:"importPath"`
	ImportName string `json:"importName"`

	ClientPkgName string `json:"clientPkgName"`
	ClientName    string `json:"clientName"`
}

func ReadConfigFile(fp string) (*Config, error) {
	b, e := ioutil.ReadFile(fp)
	if e != nil {
		return nil, e
	}
	c := new(Config)
	e = json.Unmarshal(b, c)
	if e != nil {
		return nil, e
	}
	return c, nil
}

const Template = `package crd

import (
	{{.ImportName}} "{{.ImportPath}}"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	"github.com/caicloud/dashboard-admin/pkg/errors"
	"github.com/caicloud/dashboard-admin/pkg/kubernetes"
)

const (
	CacheName{{.Name}} = "{{.Name}}"
)

func (scc *subClusterCaches) Get{{.Name}}Cache() (*{{.Plural}}Cache, bool) {
	return scc.GetAs{{.Name}}Cache(CacheName{{.Name}})
}
func (scc *subClusterCaches) GetAs{{.Name}}Cache(name string) (*{{.Plural}}Cache, bool) {
	c, ok := scc.m[name]
	if ok {
		return &{{.Plural}}Cache{lwCache: c, kc: scc.kc}, true
	}
	return nil, false
}

type {{.Plural}}Cache struct {
	lwCache *ListWatchCache
	kc      kubernetes.Interface
}

func New{{.Plural}}Cache(kc kubernetes.Interface) (*{{.Plural}}Cache, error) {
	listWatcher, objType := Get{{.Name}}CacheConfig(kc)
	c, e := NewListWatchCache(listWatcher, objType)
	if e != nil {
		return nil, e
	}
	return &{{.Plural}}Cache{
		lwCache: c,
		kc:      kc,
	}, nil
}

func (tc *{{.Plural}}Cache) Run(stopCh chan struct{}) {
	tc.lwCache.Run(stopCh)
}
{{if .IsNonNamespaced}}
func (tc *{{.Plural}}Cache) Get(key string) (*{{.ImportName}}.{{.Name}}, error) {
	return CacheGet{{.Name}}(key, tc.lwCache.indexer, tc.kc)
}
func (tc *{{.Plural}}Cache) List() ([]{{.ImportName}}.{{.Name}}, error) {
	return CacheList{{.Plural}}(tc.lwCache.indexer, tc.kc)
}
func (tc *{{.Plural}}Cache) ListCachePointer() (re []*{{.ImportName}}.{{.Name}}) {
	return CacheList{{.Plural}}Pointer(tc.lwCache.indexer, tc.kc)
}
{{else}}
func (tc *{{.Plural}}Cache) Get(namespace, key string) (*{{.ImportName}}.{{.Name}}, error) {
	return CacheGet{{.Name}}(namespace, key, tc.lwCache.indexer, tc.kc)
}
func (tc *{{.Plural}}Cache) List(namespace string) ([]{{.ImportName}}.{{.Name}}, error) {
	return CacheList{{.Plural}}(namespace, tc.lwCache.indexer, tc.kc)
}
func (tc *{{.Plural}}Cache) ListCachePointer(namespace string) (re []*{{.ImportName}}.{{.Name}}) {
	return CacheList{{.Plural}}Pointer(namespace, tc.lwCache.indexer, tc.kc)
}
{{end}}
func (tc *{{.Plural}}Cache) Indexes() cache.Indexer {
	return tc.lwCache.indexer
}

func Get{{.Name}}CacheConfig(kc kubernetes.Interface) (cache.ListerWatcher, runtime.Object) {
	return &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			options.FieldSelector = fields.Everything().String()
			return {{if .IsNonNamespaced}}kc.{{.ClientPkgName}}().{{.ClientName}}().List(options){{else}}kc.{{.ClientPkgName}}().{{.ClientName}}(metav1.NamespaceAll).List(options){{end}}
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			options.FieldSelector = fields.Everything().String()
			options.Watch = true
			return {{if .IsNonNamespaced}}kc.{{.ClientPkgName}}().{{.ClientName}}().Watch(options){{else}}kc.{{.ClientPkgName}}().{{.ClientName}}(metav1.NamespaceAll).Watch(options){{end}}
		},
	}, &{{.ImportName}}.{{.Name}}{}
}

{{if .IsNonNamespaced}}func CacheGet{{.Name}}(key string, indexer cache.Indexer, kc kubernetes.Interface) (*{{.ImportName}}.{{.Name}}, error) {
{{else}}func CacheGet{{.Name}}(namespace, key string, indexer cache.Indexer, kc kubernetes.Interface) (*{{.ImportName}}.{{.Name}}, error) {
{{end}}	if indexer != nil {
		if obj, exist, e := indexer.GetByKey(key); exist && obj != nil && e == nil {
			{{if .IsNonNamespaced}}if {{.VarName}}, _ := obj.(*{{.ImportName}}.{{.Name}}); {{.VarName}} != nil && {{.VarName}}.Name == key {
			{{else}}if {{.VarName}}, _ := obj.(*{{.ImportName}}.{{.Name}}); CheckNamespace({{.VarName}}, namespace) && {{.VarName}}.Name == key {
			{{end}}	return {{.VarName}}, nil
			}
		}
	}
	if kc == nil {
		return nil, errors.ErrVarKubeClientNil
	}
	{{.VarName}}, e := kc.{{.ClientPkgName}}().{{.ClientName}}({{if .IsNonNamespaced}}{{else}}namespace{{end}}).Get(key, metav1.GetOptions{})
	if e != nil {
		return nil, e
	}
	return {{.VarName}}, nil
}

{{if .IsNonNamespaced}}func CacheList{{.Plural}}(indexer cache.Indexer, kc kubernetes.Interface) ([]{{.ImportName}}.{{.Name}}, error) {
{{else}}func CacheList{{.Plural}}(namespace string, indexer cache.Indexer, kc kubernetes.Interface) ([]{{.ImportName}}.{{.Name}}, error) {
{{end}}	if items := indexer.List(); len(items) > 0 {
		re := make([]{{.ImportName}}.{{.Name}}, 0, len(items))
		for _, obj := range items {
			{{.VarName}}, _ := obj.(*{{.ImportName}}.{{.Name}})
			{{if .IsNonNamespaced}}if {{.VarName}} != nil {
			{{else}}if CheckNamespace({{.VarName}}, namespace) {
			{{end}}	re = append(re, *{{.VarName}})
			}
		}
		if len(re) > 0 {
			return re, nil
		}
	}
	{{.VarName}}List, e := kc.{{.ClientPkgName}}().{{.ClientName}}({{if .IsNonNamespaced}}{{else}}namespace{{end}}).List(metav1.ListOptions{})
	if e != nil {
		return nil, e
	}
	return {{.VarName}}List.Items, nil
}

{{if .IsNonNamespaced}}func CacheList{{.Plural}}Pointer(indexer cache.Indexer, kc kubernetes.Interface) (re []*{{.ImportName}}.{{.Name}}) {
{{else}}func CacheList{{.Plural}}Pointer(namespace string, indexer cache.Indexer, kc kubernetes.Interface) (re []*{{.ImportName}}.{{.Name}}) {
{{end}}	// from cache
	items := indexer.List()
	if len(items) > 0 {
		re = make([]*{{.ImportName}}.{{.Name}}, 0, len(items))
		for _, obj := range items {
			{{.VarName}}, _ := obj.(*{{.ImportName}}.{{.Name}})
			{{if .IsNonNamespaced}}if {{.VarName}} != nil {
			{{else}}if CheckNamespace({{.VarName}}, namespace) {
			{{end}}	re = append(re, {{.VarName}})
			}
		}
	}
	if len(re) > 0 {
		return re
	}
	// from source
	{{.VarName}}List, e := kc.{{.ClientPkgName}}().{{.ClientName}}({{if .IsNonNamespaced}}{{else}}namespace{{end}}).List(metav1.ListOptions{})
	if e != nil || len({{.VarName}}List.Items) == 0 {
		return nil
	}
	re = make([]*{{.ImportName}}.{{.Name}}, len({{.VarName}}List.Items))
	for i := range {{.VarName}}List.Items {
		re[i] = &{{.VarName}}List.Items[i]
	}
	return re
}
`
