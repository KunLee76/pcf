package consumer

import (
	"fmt"
	"net/http"

	"github.com/antihax/optional"

	"github.com/free5gc/openapi/Nnrf_NFDiscovery"
	"github.com/free5gc/openapi/models"
	pcf_context "github.com/free5gc/pcf/internal/context"
	"github.com/free5gc/pcf/internal/logger"
	"github.com/free5gc/pcf/internal/util"
)

func SendSearchNFInstances(
	nrfUri string, targetNfType, requestNfType models.NfType, param Nnrf_NFDiscovery.SearchNFInstancesParamOpts) (
	*models.SearchResult, error,
) {
	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath(nrfUri)
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	ctx, _, err := pcf_context.GetSelf().GetTokenCtx("nnrf-disc", "NRF")
	if err != nil {
		return nil, err
	}

	result, res, err := client.NFInstancesStoreApi.SearchNFInstances(ctx, targetNfType, requestNfType, &param)
	if err != nil {
		logger.ConsumerLog.Errorf("SearchNFInstances failed: %+v", err)
	}
	defer func() {
		if resCloseErr := res.Body.Close(); resCloseErr != nil {
			logger.ConsumerLog.Errorf("NFInstancesStoreApi response body cannot close: %+v", resCloseErr)
		}
	}()
	if res != nil && res.StatusCode == http.StatusTemporaryRedirect {
		return nil, fmt.Errorf("Temporary Redirect For Non NRF Consumer")
	}

	return &result, nil
}

func SendNFInstancesUDR(nrfUri, id string) string {
	targetNfType := models.NfType_UDR
	requestNfType := models.NfType_PCF
	localVarOptionals := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		// 	DataSet: optional.NewInterface(models.DataSetId_SUBSCRIPTION),
	}
	// switch types {
	// case NFDiscoveryToUDRParamSupi:
	// 	localVarOptionals.Supi = optional.NewString(id)
	// case NFDiscoveryToUDRParamExtGroupId:
	// 	localVarOptionals.ExternalGroupIdentity = optional.NewString(id)
	// case NFDiscoveryToUDRParamGpsi:
	// 	localVarOptionals.Gpsi = optional.NewString(id)
	// }

	result, err := SendSearchNFInstances(nrfUri, targetNfType, requestNfType, localVarOptionals)
	if err != nil {
		logger.ConsumerLog.Error(err.Error())
		return ""
	}
	for _, profile := range result.NfInstances {
		if uri := util.SearchNFServiceUri(profile, models.ServiceName_NUDR_DR, models.NfServiceStatus_REGISTERED); uri != "" {
			return uri
		}
	}
	return ""
}

func SendNFInstancesBSF(nrfUri string) string {
	targetNfType := models.NfType_BSF
	requestNfType := models.NfType_PCF
	localVarOptionals := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{}

	result, err := SendSearchNFInstances(nrfUri, targetNfType, requestNfType, localVarOptionals)
	if err != nil {
		logger.ConsumerLog.Error(err.Error())
		return ""
	}
	for _, profile := range result.NfInstances {
		if uri := util.SearchNFServiceUri(profile, models.ServiceName_NBSF_MANAGEMENT,
			models.NfServiceStatus_REGISTERED); uri != "" {
			return uri
		}
	}
	return ""
}

func SendNFInstancesAMF(nrfUri string, guami models.Guami, serviceName models.ServiceName) string {
	targetNfType := models.NfType_AMF
	requestNfType := models.NfType_PCF

	localVarOptionals := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		Guami: optional.NewInterface(util.MarshToJsonString(guami)),
	}
	// switch types {
	// case NFDiscoveryToUDRParamSupi:
	// 	localVarOptionals.Supi = optional.NewString(id)
	// case NFDiscoveryToUDRParamExtGroupId:
	// 	localVarOptionals.ExternalGroupIdentity = optional.NewString(id)
	// case NFDiscoveryToUDRParamGpsi:
	// 	localVarOptionals.Gpsi = optional.NewString(id)
	// }

	result, err := SendSearchNFInstances(nrfUri, targetNfType, requestNfType, localVarOptionals)
	if err != nil {
		logger.ConsumerLog.Error(err.Error())
		return ""
	}
	for _, profile := range result.NfInstances {
		return util.SearchNFServiceUri(profile, serviceName, models.NfServiceStatus_REGISTERED)
	}
	return ""
}
