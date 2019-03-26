package validator

// MessageMap is a map of string, that can be used as error message for ValidateStruct function.
var MessageMap = map[string]string{
	"accepted":           "The {{.Attribute}} must be accepted.",
	"activeUrl":          "The {{.Attribute}} is not a valid URL.",
	"after":              "The {{.Attribute}} must be a date after {{.Date}}.",
	"afterOrEqual":       "The {{.Attribute}} must be a date after or equal to {{.Date}}.",
	"alpha":              "The {{.Attribute}} may only contain letters.",
	"alphaDash":          "The {{.Attribute}} may only contain letters, numbers, dashes and underscores.",
	"alphaNum":           "The {{.Attribute}} may only contain letters and numbers.",
	"array":              "The {{.Attribute}} must be an array.",
	"before":             "The {{.Attribute}} must be a date before {{.Date}}.",
	"beforeOrEqual":      "The {{.Attribute}} must be a date before or equal to {{.Date}}.",
	"between.numeric":    "The {{.Attribute}} must be between {{.Min}} and {{.Max}}.",
	"between.file":       "The {{.Attribute}} must be between {{.Min}} and {{.Max}} kilobytes.",
	"between.string":     "The {{.Attribute}} must be between {{.Min}} and {{.Max}} characters.",
	"between.array":      "The {{.Attribute}} must have between {{.Min}} and {{.Max}} items.",
	"boolean":            "The {{.Attribute}} field must be true or false.",
	"confirmed":          "The {{.Attribute}} confirmation does not match.",
	"date":               "The {{.Attribute}} is not a valid date.",
	"dateFormat":         "The {{.Attribute}} does not match the format {{.Format}}.",
	"different":          "The {{.Attribute}} and {{.Other}} must be different.",
	"digits":             "The {{.Attribute}} must be {{.Digits}} digits.",
	"digitsBetween":      "The {{.Attribute}} must be between {{.Min}} and {{.Max}} digits.",
	"dimensions":         "The {{.Attribute}} has invalid image dimensions.",
	"distinct":           "The {{.Attribute}} field has a duplicate value.",
	"email":              "The {{.Attribute}} must be a valid email address.",
	"exists":             "The selected {{.Attribute}} is invalid.",
	"file":               "The {{.Attribute}} must be a file.",
	"filled":             "The {{.Attribute}} field must have a value.",
	"gt.numeric":         "The {{.Attribute}} must be greater than {{.Value}}.",
	"gt.file":            "The {{.Attribute}} must be greater than {{.Value}} kilobytes.",
	"gt.string":          "The {{.Attribute}} must be greater than {{.Value}} characters.",
	"gt.array":           "The {{.Attribute}} must have greater than {{.Value}} items.",
	"gte.numeric":        "The {{.Attribute}} must be greater than or equal {{.Value}}.",
	"gte.file":           "The {{.Attribute}} must be greater than or equal {{.Value}} kilobytes.",
	"gte.string":         "The {{.Attribute}} must be greater than or equal {{.Value}} characters.",
	"gte.array":          "The {{.Attribute}} must have {{.Value}} items or more.",
	"image":              "The {{.Attribute}} must be an image.",
	"in":                 "The selected {{.Attribute}} is invalid.",
	"inArray":            "The {{.Attribute}} field does not exist in {{.Other}}.",
	"integer":            "The {{.Attribute}} must be an integer.",
	"ip":                 "The {{.Attribute}} must be a valid IP address.",
	"ipv4":               "The {{.Attribute}} must be a valid IPv4 address.",
	"ipv6":               "The {{.Attribute}} must be a valid IPv6 address.",
	"json":               "The {{.Attribute}} must be a valid JSON string.",
	"lt.numeric":         "The {{.Attribute}} must be less than {{.Value}}.",
	"lt.file":            "The {{.Attribute}} must be less than {{.Value}} kilobytes.",
	"lt.string":          "The {{.Attribute}} must be less than {{.Value}} characters.",
	"lt.array":           "The {{.Attribute}} must have less than {{.Value}} items.",
	"lte.numeric":        "The {{.Attribute}} must be less than or equal {{.Value}}.",
	"lte.file":           "The {{.Attribute}} must be less than or equal {{.Value}} kilobytes.",
	"lte.string":         "The {{.Attribute}} must be less than or equal {{.Value}} characters.",
	"lte.array":          "The {{.Attribute}} must not have more than {{.Value}} items.",
	"max.numeric":        "The {{.Attribute}} may not be greater than {{.Max}}.",
	"max.file":           "The {{.Attribute}} may not be greater than {{.Max}} kilobytes.",
	"max.string":         "The {{.Attribute}} may not be greater than {{.Max}} characters.",
	"max.array":          "The {{.Attribute}} may not have more than {{.Max}} items.",
	"mimes":              "The {{.Attribute}} must be a file of type: {{.Values}}.",
	"mimetypes":          "The {{.Attribute}} must be a file of type: {{.Values}}.",
	"min.numeric":        "The {{.Attribute}} must be at least {{.Min}}.",
	"min.file":           "The {{.Attribute}} must be at least {{.Min}} kilobytes.",
	"min.string":         "The {{.Attribute}} must be at least {{.Min}} characters.",
	"min.array":          "The {{.Attribute}} must have at least {{.Min}} items.",
	"notIn":              "The selected {{.Attribute}} is invalid.",
	"notRegex":           "The {{.Attribute}} format is invalid.",
	"numeric":            "The {{.Attribute}} must be a number.",
	"present":            "The {{.Attribute}} field must be present.",
	"regex":              "The {{.Attribute}} format is invalid.",
	"required":           "The {{.Attribute}} field is required.",
	"requiredIf":         "The {{.Attribute}} field is required when {{.Other}} is {{.Value}}.",
	"requiredUnless":     "The {{.Attribute}} field is required unless {{.Other}} is in {{.Values}}.",
	"requiredWith":       "The {{.Attribute}} field is required when {{.Values}} is present.",
	"requiredWithAll":    "The {{.Attribute}} field is required when {{.Values}} is present.",
	"requiredWithout":    "The {{.Attribute}} field is required when {{.Values}} is not present.",
	"requiredWithoutAll": "The {{.Attribute}} field is required when none of {{.Values}} are present.",
	"same":               "The {{.Attribute}} and {{.Other}} must match.",
	"size.numeric":       "The {{.Attribute}} must be {{.Size}}.",
	"size.file":          "The {{.Attribute}} must be {{.Size}} kilobytes.",
	"size.string":        "The {{.Attribute}} must be {{.Size}} characters.",
	"size.array":         "The {{.Attribute}} must contain {{.Size}} items.",
	"string":             "The {{.Attribute}} must be a string.",
	"timezone":           "The {{.Attribute}} must be a valid zone.",
	"unique":             "The {{.Attribute}} has already been taken.",
	"uploaded":           "The {{.Attribute}} failed to upload.",
	"url":                "The {{.Attribute}} format is invalid.",
}