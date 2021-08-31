package plugin

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	// "github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/valyala/fastjson"
)

var (
	_ backend.QueryDataHandler      = (*SampleDatasource)(nil)
	_ backend.CheckHealthHandler    = (*SampleDatasource)(nil)
	_ backend.CallResourceHandler   = (*SampleDatasource)(nil)
	_ instancemgmt.InstanceDisposer = (*SampleDatasource)(nil)
)

type SampleDatasource struct{}

func NewSampleDatasource(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	return &SampleDatasource{}, nil
}

func (d *SampleDatasource) Dispose() {}

func makeHttpClient(settings *backend.DataSourceInstanceSettings) (*http.Client, error) {
	if settings == nil {
		return nil, errors.New("no http settings for this request")
	}

	clientOptions, err := settings.HTTPClientOptions()
	if err != nil {
		return nil, err
	}
	client, err := httpclient.New(clientOptions)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (d *SampleDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	response := backend.NewQueryDataResponse()
	for _, q := range req.Queries {
		response.Responses[q.RefID] = d.doDataQuery(ctx, req.PluginContext.DataSourceInstanceSettings, q)
	}
	return response, nil
}

func (d *SampleDatasource) doDataQuery(ctx context.Context, settings *backend.DataSourceInstanceSettings, q backend.DataQuery) backend.DataResponse {
	var start int64
	var step int64
	var end int64
	var maxRows int64

	start = q.TimeRange.From.Unix()
	end = q.TimeRange.To.Unix()
	step = int64(q.Interval.Seconds())
	if step == 0 {
		step = 1
	}
	maxRows = q.MaxDataPoints
	if maxRows < 10 { // RRDTool doesn't accept less than 10.
		maxRows = 10
	}

	var p fastjson.Parser
	queryJson, err := p.ParseBytes(q.JSON)
	if err != nil {
		return backend.DataResponse{
			Error: fmt.Errorf("unable to parse query json fields: %s", err),
		}
	}
	payload := string(queryJson.GetStringBytes("payload"))

	queryPath := fmt.Sprintf(
		"/api/v1/xport?start=%d&step=%d&end=%d&maxrows=%d&xport=%s",
		start,
		step,
		end,
		maxRows,
		url.QueryEscape(payload),
	)

	responseBytes, err := d.doGet(ctx, settings, queryPath)
	if err != nil {
		return backend.DataResponse{
			Error: err,
		}
	}

	response, err := p.ParseBytes(responseBytes)
	if err != nil {
		return backend.DataResponse{
			Error: fmt.Errorf("unable to parse response json: %s", err),
		}
	}

	meta := response.Get("meta")
	legend := meta.GetArray("legend")
	dataStart, dataStartErr := meta.Get("start").Int64()
	dataStep, dataStepErr := meta.Get("step").Int64()
	dataRows := response.GetArray("data")

	if meta == nil || dataRows == nil || legend == nil || dataStartErr != nil || dataStepErr != nil {
		return backend.DataResponse{
			Error: errors.New("response missing required fields"),
		}
	}

	frames := []*data.Frame{}

	for i, nameJsonValue := range legend {
		name := string(nameJsonValue.GetStringBytes())
		if name == "" {
			name = fmt.Sprintf("Value%d", i)
		}

		timeCol := make([]time.Time, 0, len(dataRows))
		dataCol := make([]float64, 0, len(dataRows))

		t := dataStart
		for _, rowJsonValue := range dataRows {
			rowData := rowJsonValue.GetArray()
			if len(rowData) != len(legend) {
				return backend.DataResponse{
					Error: errors.New("data had incorrect number of colums"),
				}
			}
			v := rowData[i]
			if v.Type() == fastjson.TypeNumber {
				timeCol = append(timeCol, time.Unix(t, 0))
				dataCol = append(dataCol, v.GetFloat64())
			}
			t += dataStep
		}
		frame := data.NewFrame(
			fmt.Sprintf("%s", name),
			data.NewField("Time", nil, timeCol),
			data.NewField(name, nil, dataCol),
		)
		frame.RefID = q.RefID
		frame.SetMeta(&data.FrameMeta{
			ExecutedQueryString: queryPath,
		})
		frames = append(frames, frame)
	}

	return backend.DataResponse{
		Frames: frames,
	}
}

func (d *SampleDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {

	if req.PluginContext.DataSourceInstanceSettings.URL == "" {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "url setting is empty",
		}, nil
	}

	body, err := d.doGet(ctx, req.PluginContext.DataSourceInstanceSettings, "/api/v1/ping")
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
	}

	if !bytes.Equal(body, []byte("\"pong\"")) {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: fmt.Sprintf("expected \"pong\", got %q ", string(body)),
		}, nil
	}

	return &backend.CheckHealthResult{
		Status: backend.HealthStatusOk,
	}, nil
}

func (d *SampleDatasource) doGet(ctx context.Context, settings *backend.DataSourceInstanceSettings, path string) ([]byte, error) {

	client, err := makeHttpClient(settings)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, "GET", settings.URL+path, nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error fetching %s: %s", path, err.Error())
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading request body: %s", err.Error())
	}

	if response.StatusCode != 200 {
		errMsg := string(body)
		if errMsg == "" {
			errMsg = response.Status
		}
		return nil, errors.New(errMsg)
	}

	return body, nil
}

func (d *SampleDatasource) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {

	settings := req.PluginContext.DataSourceInstanceSettings

	client, err := makeHttpClient(settings)
	if err != nil {
		return err
	}

	var bodyReader io.Reader = nil

	if req.Body != nil {
		bodyReader = bytes.NewReader(req.Body)
	} 

	request, err := http.NewRequestWithContext(ctx, req.Method, settings.URL+req.URL, bodyReader)
	if err != nil {
		return err
	}
	res, err := client.Do(request)
	if err != nil {
		return sender.Send(&backend.CallResourceResponse{
			Status:  502,
			Headers: res.Header,
			Body:    []byte(fmt.Sprintf("error performing request: %s", err)),
		})
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return sender.Send(&backend.CallResourceResponse{
			Status:  502,
			Headers: res.Header,
			Body:    []byte(fmt.Sprintf("error reading response: %s", err)),
		})
	}

	return sender.Send(&backend.CallResourceResponse{
		Status:  res.StatusCode,
		Headers: res.Header,
		Body:    body,
	})
}
