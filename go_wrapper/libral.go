package libral

/*
#cgo LDFLAGS: -fPIC -L${SRCDIR}/ -lstdc++ -lral -lcpp-hocon -lyaml-cpp -l:leatherman_execution.a -l:leatherman_logging.a -l:leatherman_locale.a -l:leatherman_ruby.a -l:leatherman_dynamic_library.a -l:leatherman_util.a -l:leatherman_file_util.a -l:leatherman_json_container.a -lboost_system-mt -lboost_log-mt -lboost_log_setup-mt -lboost_thread-mt -lboost_date_time -lboost_filesystem-mt -lboost_chrono-mt -lboost_regex-mt -lboost_atomic-mt -lboost_program_options-mt -lboost_locale-mt -laugeas -liconv -lselinux -lfa -lstdc++ -lm
#cgo CFLAGS: -I${SRCDIR}/../lib/inc

#include "libral/cwrapper.hpp"
#include <stdlib.h>
*/
import "C"

import (
	"encoding/json"
	"fmt"
	"libral/types"
	"unsafe"
)

// GetProviders returns the resource providers available to libral.
func GetProviders() ([]types.Provider, error) {
	rawProviders, err := getProvidersRaw()
	if err != nil {
		return nil, err
	}

	var providersResult types.ProvidersResult

	if err := json.Unmarshal([]byte(rawProviders), &providersResult); err != nil {
		return nil, fmt.Errorf("Error unmarshalling provider JSON: %v", err)
	}

	return providersResult.Providers, nil
}

// GetResources returns all resources of the specified provider type.
func GetResources(typeName string) ([]types.Resource, error) {
	var result []types.Resource
	rawResources, err := getResourcesRaw(typeName)
	if err != nil {
		return nil, err
	}

	var resourcesResult types.ResourcesResult

	if err := json.Unmarshal([]byte(rawResources), &resourcesResult); err != nil {
		return nil, fmt.Errorf("Error unmarshalling resources JSON: %v", err)
	}

	for _, rawResource := range resourcesResult.Resources {
		var resource types.Resource
		if err := json.Unmarshal([]byte(rawResource), &resource); err != nil {
			return nil, fmt.Errorf("Error unmarshalling resource JSON: %v", err)
		}
		resource.Raw = rawResource
		var attributes map[string]interface{}
		if err := json.Unmarshal([]byte(rawResource), &attributes); err != nil {
			return nil, fmt.Errorf("Error unmarshalling resource attributes JSON: %v", err)
		}
		_, ok := attributes["ral"]
		if ok {
			delete(attributes, "ral")
		}
		resource.Attributes = attributes
		result = append(result, resource)
	}

	return result, nil
}

// GetResource returns a Resource of the specified provider type with matching resource name,
// or an empty resource if no match is found, if more than one resource exists with the same
// name an error is thrown.
func GetResource(typeName, resourceName string) (types.Resource, error) {
	var result types.Resource
	rawResource, err := getResourceRaw(typeName, resourceName)
	if err != nil {
		return types.Resource{}, err
	}

	var resourceResult types.ResourceResult
	if err := json.Unmarshal([]byte(rawResource), &resourceResult); err != nil {
		return types.Resource{}, fmt.Errorf("Error unmarshalling resource JSON: %v", err)
	}

	if err := json.Unmarshal([]byte(resourceResult.Resource), &result); err != nil {
		return types.Resource{}, fmt.Errorf("Error unmarshalling resource details JSON: %v", err)
	}
	result.Raw = resourceResult.Resource
	var attributes map[string]interface{}
	if err := json.Unmarshal([]byte(resourceResult.Resource), &attributes); err != nil {
		return types.Resource{}, fmt.Errorf("Error unmarshalling resource attributes JSON: %v", err)
	}
	_, ok := attributes["ral"]
	if ok {
		delete(attributes, "ral")
	}
	result.Attributes = attributes

	return result, nil
}

// getProvidersRaw returns a JSON string containing a list of providers known to libral.
//
// For example:
//
//   {
//       "providers": [
//           {
//               "name": "file::posix",
//               "type": "file",
//               "source": "builtin",
//               "desc": "A provider to manage POSIX files.\n",
//               "suitable": true,
//               "attributes": [
//                   {
//                       "name": "checksum",
//                       "desc": "(missing description)",
//                       "type": "enum[md5, md5lite, sha256, sha256lite, mtime, ctime, none]",
//                       "kind": "r"
//                   },
//                   ...
//               ]
//           },
//           {
//               "name": "host::aug",
//               "type": "host",
//               ...
//           }
//       ]
//   }
func getProvidersRaw() (string, error) {
	var resultC *C.char
	defer C.free(unsafe.Pointer(resultC))

	ok := C.get_providers(&resultC)
	if ok != 0 {
		return "", fmt.Errorf("Error thrown calling get_providers: %d", ok)
	}
	result := C.GoString(resultC)
	return result, nil
}

// getResourcesRaw returns a JSON string containing a list all resources of the specified type
// found by libral.
//
// For example, querying the `host` type:
//
//   {
//       "resources": [
//           {
//               "name": "localhost",
//               "ensure": "present",
//               "ip": "127.0.0.1",
//               "target": "/etc/hosts",
//               "ral": {
//                   "type": "host",
//                   "provider": "host::aug"
//               }
//           },
//           {
//               "name": "broadcasthost",
//               "ensure": "present",
//               "ip": "255.255.255.255",
//               "target": "/etc/hosts",
//               "ral": {
//                   "type": "host",
//                   "provider": "host::aug"
//               }
//           }
//       ]
//   }
func getResourcesRaw(typeName string) (string, error) {
	var resultC *C.char
	defer C.free(unsafe.Pointer(resultC))
	typeNameC := C.CString(typeName)
	defer C.free(unsafe.Pointer(typeNameC))

	ok := C.get_resources(&resultC, typeNameC)
	if ok != 0 {
		return "", fmt.Errorf("Error thrown calling get_all_resources: %d", ok)
	}
	result := C.GoString(resultC)
	return result, nil
}

// getResourceRaw returns a JSON string containing the matching resources of the specified type
// found by libral. It will throw an error if multiple resources are found with the same name.
//
// For example, querying the `host` type `broadcasthost`:
//
//   {
//       "resources": [
//           {
//               "name": "broadcasthost",
//               "ensure": "present",
//               "ip": "255.255.255.255",
//               "target": "/etc/hosts",
//               "ral": {
//                   "type": "host",
//                   "provider": "host::aug"
//               }
//           }
//       ]
//   }
func getResourceRaw(typeName, resourceName string) (string, error) {
	var resultC *C.char
	defer C.free(unsafe.Pointer(resultC))
	typeNameC := C.CString(typeName)
	defer C.free(unsafe.Pointer(typeNameC))
	resourceNameC := C.CString(resourceName)
	defer C.free(unsafe.Pointer(resourceNameC))

	ok := C.get_resource(&resultC, typeNameC, resourceNameC)
	if ok != 0 {
		return "", fmt.Errorf("Error thrown calling get_all_resources: %d", ok)
	}
	result := C.GoString(resultC)
	return result, nil
}
