package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/spf13/pflag"
	Digital_Ocean_Clusterv1alpha1 "github.com/anutosh491/Digital_Ocean_Cluster_Operator/pkg/apis/anutosh491.dev/v1alpha1"
	admv1beta1 "k8s.io/api/admission/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/options"
	"k8s.io/component-base/cli/globalflag"

	kdo "github.com/anutosh491/Valid_Controller/pkg/digitalocean"
)

type Options struct {
	SecureServingOptions options.SecureServingOptions
}

func (o *Options) AddFlagSet(fs *pflag.FlagSet) {
	o.SecureServingOptions.AddFlags(fs)
}

type Config struct {
	SecureServingInfo *server.SecureServingInfo
}

func (o *Options) Config() *Config {
	if err := o.SecureServingOptions.MaybeDefaultWithSelfSignedCerts("0.0.0.0", nil, nil); err != nil {
		panic(err)
	}

	c := Config{}
	o.SecureServingOptions.ApplyTo(&c.SecureServingInfo)
	return &c
}

const (
	valKon = "Valid-Controller"
)

func NewDefaultOptions() *Options {
	o := &Options{
		SecureServingOptions: *options.NewSecureServingOptions(),
	}
	o.SecureServingOptions.BindPort = 8443
	o.SecureServingOptions.ServerCert.PairName = valKon
	return o
}

func main() {
	options := NewDefaultOptions()

	fs := pflag.NewFlagSet(valKon, pflag.ExitOnError)
	globalflag.AddGlobalFlags(fs, valKon)

	options.AddFlagSet(fs)

	if err := fs.Parse(os.Args); err != nil {
		panic(err)
	}

	c := options.Config()

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(ServeDigital_Ocean_ClusterValidation))

	stopCh := server.SetupSignalHandler()
	ch, err := c.SecureServingInfo.Serve(mux, 30*time.Second, stopCh)
	if err != nil {
		panic(err)
	} else {
		<-ch
	}
}

var (
	scheme = runtime.NewScheme()
	codecs = serializer.NewCodecFactory(scheme)
)

func ServeDigital_Ocean_ClusterValidation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ServeDigital_Ocean_ClusterValidation was called")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responsewriters.InternalError(w, r, err)
		fmt.Printf("Error %s, reading the body", err.Error())
	}

	// read body to get an instance of admissionreivew object
	// get the gvk for admission review
	gvk := admv1beta1.SchemeGroupVersion.WithKind("AdmissionReview")
	// var of type admission reiveew
	var admissionReview admv1beta1.AdmissionReview
	_, _, err = codecs.UniversalDeserializer().Decode(body, &gvk, &admissionReview)
	if err != nil {
		fmt.Printf("Error %s, converting req body to admission review type", err.Error())

	}

	// get Digital_Ocean_Cluster spec from admission review object
	gvkDigital_Ocean_Cluster := Digital_Ocean_Clusterv1alpha1.SchemeGroupVersion.WithKind("Digital_Ocean_Cluster")
	var k Digital_Ocean_Clusterv1alpha1.Digital_Ocean_Cluster
	_, _, err = codecs.UniversalDeserializer().Decode(admissionReview.Request.Object.Raw, &gvkDigital_Ocean_Cluster, &k)
	if err != nil {
		fmt.Printf("Error %s, while getting Digital_Ocean_Cluster type from admission review", err.Error())
	}

	response := admv1beta1.AdmissionResponse{}
	allow := validateDigital_Ocean_Cluster(k.Spec)
	// if allow is false or err is not nil
	if !allow {
		response = admv1beta1.AdmissionResponse{
			UID:     admissionReview.Request.UID,
			Allowed: allow,
			Result: &v1.Status{
				Message: fmt.Sprintf("The specified version %s is not supported by DO", k.Spec.Version),
			},
		}
	} else {
		response = admv1beta1.AdmissionResponse{
			UID:     admissionReview.Request.UID,
			Allowed: allow,
		}
	}

	admissionReview.Response = &response
	// write the response to response writer
	fmt.Printf("respoknse that we are trying to return is %+v\n", response)
	res, err := json.Marshal(admissionReview)
	if err != nil {
		fmt.Printf("error %s, while converting response to byte slice", err.Error())
	}

	_, err = w.Write(res)
	if err != nil {
		fmt.Printf("error %s, writing respnse to responsewriter", err.Error())
	}
}

func validateDigital_Ocean_Cluster(spek Digital_Ocean_Clusterv1alpha1.Digital_Ocean_ClusterSpec) bool {
	_, err := kdo.ValidateDigital_Ocean_ClusterVersion(spek)
	if err != nil {
		fmt.Printf("error %s vaidating Digital_Ocean_Cluster resource ", err.Error())
		return false
	}

	return true
}
