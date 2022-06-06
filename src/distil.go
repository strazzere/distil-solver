// RedNaga / Tim Strazzere (c) 2018-*

package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"github.com/robertkrimen/otto"
)

// ceiling = Math.pow(2, 32 - iterations)
// for (var i = 0;;) {
// 	var o = i.toString(16) + ":" + blobString;
// 	i++;
// 	var hashOutput = sha1(o);
// 	if (parseInt(hashOutput.substr(0, 8), 16) < ceiling) {
// 		return void callback(o)
// 	}
// }
func workOnProof(blob string, iterations int) (string, error) {
	ceiling := int64(math.Pow(2, float64(32-iterations)))
	i := 0
	for {
		hashInput := fmt.Sprintf("%x:%s", i, blob)
		hash := sha1.Sum([]byte(hashInput))
		val, err := strconv.ParseInt(fmt.Sprintf("%x", hash[0:4]), 16, 64)
		if err != nil {
			log.Printf("Error : %v", err)
			return "", err
		}
		if val < ceiling {
			return hashInput, nil
		}
		i++
	}
}

func getProof() string {
	return fmt.Sprintf("%d:%s", time.Now().UnixNano()/int64(time.Millisecond), setUpPow(20))
}

// function t(e) {
// 	for (var t = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", n = "", r = 0; e > r; ++r) {
// 		n += t.substr(Math.floor(Math.random() * t.length), 1);
// 	}
// 	return n
// }
// 	r = (new Date).getTime() + ":" + t(20);
func setUpPow(iterations int) string {
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	out := ""
	var template = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	for i := 0; iterations > i; i++ {
		charStart := int(math.Floor(random.Float64() * float64(len(template))))
		out = fmt.Sprintf("%s%s", out, template[charStart:charStart+1])
	}

	return out
}

func instrumentWindow(vm *otto.Otto) error {
	window, err := vm.Object(`window = {}`)
	if err != nil {
		return err
	}
	addEventListener := func(call otto.FunctionCall) otto.Value {
		fmt.Printf("addEventListener called via window : %+v\n", call)
		//		_, value := newTimer(call, true)
		return otto.Value{}
	}
	window.Set("addEventListener", addEventListener)

	attachEvent := func(call otto.FunctionCall) otto.Value {
		// fmt.Printf("attachEvent called : %+v\n", call)
		//		_, value := newTimer(call, true)
		return otto.Value{}
	}
	window.Set("attachEvent", attachEvent)

	err = instrumentDocument(vm)
	if err != nil {
		return err
	}
	err = instrumentXMLHttpRequest(vm)
	if err != nil {
		return err
	}
	window.Set("XMLHttpRequest", true)

	return nil
}

func instrumentDocument(vm *otto.Otto) error {
	document, err := vm.Object(`document = {}`)
	if err != nil {
		return err
	}
	window, err := vm.Object(`window`)
	if err != nil {
		return err
	}
	window.Set("document", document)

	getElementByID := func(call otto.FunctionCall) otto.Value {
		// fmt.Printf("getElementById called : %+v\n", call)
		//		_, value := newTimer(call, true)
		return otto.Value{}
	}
	document.Set("getElementById", getElementByID)

	addEventListener := func(call otto.FunctionCall) otto.Value {
		fmt.Printf("addEventListener called via document : %+v\n", call)
		//		_, value := newTimer(call, true)
		return otto.Value{}
	}
	document.Set("addEventListener", addEventListener)

	document.Set("readyState", "complete")

	document.Set("audio", "loading")

	return nil
}

func instrumentXMLHttpRequest(vm *otto.Otto) error {
	XMLHttpRequest, err := vm.Object(`XMLHttpRequest = {}`)
	if err != nil {
		return err
	}

	open := func(call otto.FunctionCall) otto.Value {
		// fmt.Printf("XMLHttpRequest.open() called : %+v\n", call)
		//		_, value := newTimer(call, true)
		return otto.Value{}
	}
	XMLHttpRequest.Set("open", open)

	send := func(call otto.FunctionCall) otto.Value {
		// fmt.Printf("XMLHttpRequest.send() called : %+v\n", call)
		//		_, value := newTimer(call, true)
		return otto.Value{}
	}
	XMLHttpRequest.Set("send", send)

	vm.Set("XMLHttpRequest", XMLHttpRequest)

	// XMLHttpRequest := func(call otto.FunctionCall) otto.Value {
	// 	fmt.Printf("XMLHttpRequest called : %+v\n", call)
	// 	//		_, value := newTimer(call, true)
	// 	return otto.Value{}
	// }
	// vm.Set("XMLHttpRequest", XMLHttpRequest)

	return nil
}

func getEscapedProofQuery(proof string, userAgent string) (string, error) {
	queryStruct, err := getProofQuery(proof, userAgent)
	if err != nil {
		return "", err
	}
	data, err := json.Marshal(queryStruct)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("p=%s", url.QueryEscape(string(data))), nil
}

func getProofQuery(proof string, userAgent string) (interface{}, error) {
	var preformedStruct interface{}

	err := json.Unmarshal([]byte(preformed), &preformedStruct)
	if err != nil {
		return "", err
	}

	preformedStruct.(map[string]interface{})["proof"] = proof

	if userAgent != "" {
		fp2 := preformedStruct.(map[string]interface{})["fp2"]
		fp2.(map[string]interface{})["userAgent"] = userAgent
		preformedStruct.(map[string]interface{})["fp2"] = fp2
	}

	return preformedStruct, nil
}

const preformed = `{
	"proof": "63:1554223568953:cqjyfXPuFws8N8xNILhN",
	"fp2": {
	  "userAgent": "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.86 Safari/537.36",
	  "language": "en-US",
	  "screen": {
		"width": 1680,
		"height": 1050,
		"availHeight": 944,
		"availWidth": 1680,
		"pixelDepth": 24,
		"innerWidth": 1379,
		"innerHeight": 352,
		"outerWidth": 1379,
		"outerHeight": 794,
		"devicePixelRatio": 2
	  },
	  "timezone": -7,
	  "indexedDb": true,
	  "addBehavior": false,
	  "openDatabase": true,
	  "cpuClass": "unknown",
	  "platform": "MacIntel",
	  "doNotTrack": "unknown",
	  "plugins": "",
	  "canvas": {
		"winding": "yes",
		"towebp": true,
		"blending": true,
		"img": "a98d9edca03b1d0f259d09b5baa73ba2844d1d14"
	  },
	  "webGL": {
		"img": "0902675f2196d7d89ea7751211da3a7db20e7401",
		"extensions": "ANGLE_instanced_arrays;EXT_blend_minmax;EXT_color_buffer_half_float;EXT_disjoint_timer_query;EXT_frag_depth;EXT_shader_texture_lod;EXT_texture_filter_anisotropic;WEBKIT_EXT_texture_filter_anisotropic;EXT_sRGB;OES_element_index_uint;OES_standard_derivatives;OES_texture_float;OES_texture_float_linear;OES_texture_half_float;OES_texture_half_float_linear;OES_vertex_array_object;WEBGL_color_buffer_float;WEBGL_compressed_texture_s3tc;WEBKIT_WEBGL_compressed_texture_s3tc;WEBGL_compressed_texture_s3tc_srgb;WEBGL_debug_renderer_info;WEBGL_debug_shaders;WEBGL_depth_texture;WEBKIT_WEBGL_depth_texture;WEBGL_draw_buffers;WEBGL_lose_context;WEBKIT_WEBGL_lose_context",
		"aliasedlinewidthrange": "[1,1]",
		"aliasedpointsizerange": "[1,255.875]",
		"alphabits": 8,
		"antialiasing": "yes",
		"bluebits": 8,
		"depthbits": 24,
		"greenbits": 8,
		"maxanisotropy": 16,
		"maxcombinedtextureimageunits": 80,
		"maxcubemaptexturesize": 16384,
		"maxfragmentuniformvectors": 1024,
		"maxrenderbuffersize": 16384,
		"maxtextureimageunits": 16,
		"maxtexturesize": 16384,
		"maxvaryingvectors": 15,
		"maxvertexattribs": 16,
		"maxvertextextureimageunits": 16,
		"maxvertexuniformvectors": 1024,
		"maxviewportdims": "[16384,16384]",
		"redbits": 8,
		"renderer": "WebKitWebGL",
		"shadinglanguageversion": "WebGLGLSLES1.0(OpenGLESGLSLES1.0Chromium)",
		"stencilbits": 0,
		"vendor": "WebKit",
		"version": "WebGL1.0(OpenGLES2.0Chromium)",
		"vertexshaderhighfloatprecision": 23,
		"vertexshaderhighfloatprecisionrangeMin": 127,
		"vertexshaderhighfloatprecisionrangeMax": 127,
		"vertexshadermediumfloatprecision": 23,
		"vertexshadermediumfloatprecisionrangeMin": 127,
		"vertexshadermediumfloatprecisionrangeMax": 127,
		"vertexshaderlowfloatprecision": 23,
		"vertexshaderlowfloatprecisionrangeMin": 127,
		"vertexshaderlowfloatprecisionrangeMax": 127,
		"fragmentshaderhighfloatprecision": 23,
		"fragmentshaderhighfloatprecisionrangeMin": 127,
		"fragmentshaderhighfloatprecisionrangeMax": 127,
		"fragmentshadermediumfloatprecision": 23,
		"fragmentshadermediumfloatprecisionrangeMin": 127,
		"fragmentshadermediumfloatprecisionrangeMax": 127,
		"fragmentshaderlowfloatprecision": 23,
		"fragmentshaderlowfloatprecisionrangeMin": 127,
		"fragmentshaderlowfloatprecisionrangeMax": 127,
		"vertexshaderhighintprecision": 0,
		"vertexshaderhighintprecisionrangeMin": 31,
		"vertexshaderhighintprecisionrangeMax": 30,
		"vertexshadermediumintprecision": 0,
		"vertexshadermediumintprecisionrangeMin": 31,
		"vertexshadermediumintprecisionrangeMax": 30,
		"vertexshaderlowintprecision": 0,
		"vertexshaderlowintprecisionrangeMin": 31,
		"vertexshaderlowintprecisionrangeMax": 30,
		"fragmentshaderhighintprecision": 0,
		"fragmentshaderhighintprecisionrangeMin": 31,
		"fragmentshaderhighintprecisionrangeMax": 30,
		"fragmentshadermediumintprecision": 0,
		"fragmentshadermediumintprecisionrangeMin": 31,
		"fragmentshadermediumintprecisionrangeMax": 30,
		"fragmentshaderlowintprecision": 0,
		"fragmentshaderlowintprecisionrangeMin": 31,
		"fragmentshaderlowintprecisionrangeMax": 30
	  },
	  "touch": {
		"maxTouchPoints": 0,
		"touchEvent": false,
		"touchStart": false
	  },
	  "video": {
		"ogg": "probably",
		"h264": "probably",
		"webm": "probably"
	  },
	  "audio": {
		"ogg": "probably",
		"mp3": "probably",
		"wav": "probably",
		"m4a": "maybe"
	  },
	  "vendor": "GoogleInc.",
	  "product": "Gecko",
	  "productSub": "20030107",
	  "browser": {
		"ie": false,
		"chrome": true,
		"webdriver": false
	  },
	  "window": {
		"historyLength": 2,
		"hardwareConcurrency": 4,
		"iframe": false
	  },
	  "fonts": "ArialUnicodeMS;GillSans;HelveticaNeue"
	},
	"cookies": 1,
	"setTimeout": 0,
	"setInterval": 0,
	"appName": "Netscape",
	"platform": "MacIntel",
	"syslang": "en-US",
	"userlang": "en-US",
	"cpu": "",
	"productSub": "20030107",
	"plugins": {},
	"mimeTypes": {},
	"screen": {
	  "width": 1680,
	  "height": 1050,
	  "colorDepth": 24
	},
	"fonts": {
	  "0": "HoeflerText",
	  "1": "Monaco",
	  "2": "Georgia",
	  "3": "TrebuchetMS",
	  "4": "Verdana",
	  "5": "AndaleMono",
	  "6": "Monaco",
	  "7": "CourierNew",
	  "8": "Courier"
	}
  }  
`
