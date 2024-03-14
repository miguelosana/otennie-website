package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"otennie/handlers"
	"otennie/storage"
	"strings"

	"github.com/spf13/viper"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	viper.SetDefault("http-port", 80)
	viper.SetDefault("use-https", false)
	viper.SetEnvPrefix("OTENNIE")
	viper.SetDefault("files-path", "../../dist")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	if _, fileErr := os.Stat("./config.yaml"); fileErr == nil {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		if readErr := viper.ReadInConfig(); readErr != nil {

			log.Fatalf("fatal error config file: %w", readErr)
		}
	}

	useHttps := viper.GetBool("use-https")
	ctx := context.Background()
	var db handlers.ModelWriter
	if viper.GetString("db-file") != "" {

		db = storage.NewBoltStorage(viper.GetString("db-file"))
	} else {
		db = storage.NewFirestoreStorage(ctx, viper.GetString("PROJECT_ID"))
	}

	s := handlers.NewServer(db, viper.GetString("files-path"))
	defer s.Close()

	var httpsSrv *http.Server
	var m *autocert.Manager
	if useHttps {
		dataDir := "."
		hostPolicy := func(ctx context.Context, host string) error {
			allowedHost := "otennie.com"
			if host == allowedHost {
				return nil
			}
			return fmt.Errorf("acme/autocert: only %s host is allowed", allowedHost)
		}

		httpsSrv = s.MakeHttpServer()
		m = &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: hostPolicy,
			Cache:      autocert.DirCache(dataDir),
		}

		httpsSrv.Addr = ":443"
		httpsSrv.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}

		go func() {
			err := httpsSrv.ListenAndServeTLS("", "")
			if err != nil {
				log.Fatalf("httpSrv.ListenAndServeTLs() failed with %v", err)
			}
		}()
	}

	var httpSrv *http.Server
	if useHttps {
		httpSrv = handlers.MakeHTTPToHTTPSRedirectServer()
	} else {
		httpSrv = s.MakeHttpServer()
	}
	if m != nil {
		httpSrv.Handler = m.HTTPHandler(httpSrv.Handler)

	}
	httpSrv.Addr = fmt.Sprintf(":%d", viper.GetInt("http-port"))
	err := httpSrv.ListenAndServe()

	if err != nil {
		log.Fatalf("httpSrv.ListenAndServe() failed with %v", err)
	}
}
