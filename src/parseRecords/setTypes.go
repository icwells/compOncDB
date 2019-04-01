// Regular expression dictionaries for the matcher struct

package main

import (
	"regexp"
)

type diagnosis struct {
	expression *regexp.Regexp
	malignant  string
}

func (m *matcher) setTypes() {
	// Sets type and location maps
	m.location = map[string]*regexp.Regexp{
		"abdomen":         regexp.MustCompile(`(?i)abdomen|abdom.*|omentum|diaphragm`),
		"bile duct":       regexp.MustCompile(`(?i)bile.*|biliary`),
		"bone":            regexp.MustCompile(`(?i)sacrum|bone.*`),
		"brain":           regexp.MustCompile(`(?i)brain`),
		"adrenal":         regexp.MustCompile(`(?i)adrenal`),
		"bladder":         regexp.MustCompile(`(?i)bladder`),
		"breast":          regexp.MustCompile(`(?i)breast|mammary`),
		"colon":           regexp.MustCompile(`(?i)colon|rectum`),
		"duodenum":        regexp.MustCompile(`(?i)duodenum`),
		"fat":             regexp.MustCompile(`(?i)fat|adipose.*`),
		"heart":           regexp.MustCompile(`(?i)heart|cardiac|atrial`),
		"small intestine": regexp.MustCompile(`(?i)(small )?intestin(e|al)`),
		"kidney":          regexp.MustCompile(`(?i)kidney.*|ureter|renal`),
		"leukemia":        regexp.MustCompile(`(?i)leukemia`),
		"liver":           regexp.MustCompile(`(?i)hepa.*|liver.*|hep.*|billia.*`),
		"lung":            regexp.MustCompile(`(?i)lung.*|pulm.*|mediasti.*|bronchial|alveol.*`),
		"lymph nodes":     regexp.MustCompile(`(?i)lymph( node)?`),
		"muscle":          regexp.MustCompile(`(?i)muscle|.*structure.*`),
		"nerve":           regexp.MustCompile(`(?i)nerve.*`),
		"other":           regexp.MustCompile(`(?i)gland|basal.*|islet|multifocal|neck|nasal|neuroendo.*`),
		"oral":            regexp.MustCompile(`(?i)oral|tongue|mouth|lip|palate|pharyn.*|laryn.*`),
		"ovary":           regexp.MustCompile(`(?i)ovar.*`),
		"pancreas":        regexp.MustCompile(`(?i)pancreas.*|islet`),
		"seminal vesicle": regexp.MustCompile(`(?i)seminal vesicle`),
		"skin":            regexp.MustCompile(`(?i)skin|eyelid|(sub)?cutan.*|derm.*`),
		"spinal cord":     regexp.MustCompile(`(?i)spinal|spine`),
		"spleen":          regexp.MustCompile(`(?i)spleen`),
		"testis":          regexp.MustCompile(`(?i)testi.*`),
		"thyroid":         regexp.MustCompile(`(?i)thyroid`),
		"uterus":          regexp.MustCompile(`(?i)uter.*`),
		"vulva":           regexp.MustCompile(`(?i)vulva|vagina`),
		"widespread":      regexp.MustCompile(`(?i)widespread|metastatic|body as a whole|multiple|disseminated`),
	}
	m.types = make(map[string]map[string]diagnosis)
	m.types["adenoma"] = map[string]diagnosis{
		"adenoma":     {regexp.MustCompile(`(?i)adenoma`), "0"},
		"cystadenoma": {regexp.MustCompile(`(?i)cystadenoma`), "0"},
	}
	m.types["carcinoma"] = map[string]diagnosis{
		"adenocarcinoma":      {regexp.MustCompile(`(?i)adenocarcinoma`), "1"},
		"carcinoma":           {regexp.MustCompile(`(?i)carcinoma|TCC`), "1"},
		"carcinomatosis":      {regexp.MustCompile(`(?i)carinomatosis`), "1"},
		"fibroadenocarcinoma": {regexp.MustCompile(`(?i)fibroadenocarcinoma`), "1"},
	}
	m.types["sarcoma"] = map[string]diagnosis{
		"sarcoma":          {regexp.MustCompile(`(?i)sarcoma`), "1"},
		"lymphoma":         {regexp.MustCompile(`(?i)lymphoma|lymphosarcoma`), "1"},
		"chondrosarcoma":   {regexp.MustCompile(`(?i)chondrosarcoma`), "1"},
		"fibrosarcoma":     {regexp.MustCompile(`(?i)fibrosarcoma`), "1"},
		"hemangiosarcoma":  {regexp.MustCompile(`(?i)Hemangiosarcoma`), "1"},
		"leiomyosarcoma":   {regexp.MustCompile(`(?i)leiomyosarcoma`), "1"},
		"myxosarcoma":      {regexp.MustCompile(`(?i)myxosarcoma`), "1"},
		"osteosarcoma":     {regexp.MustCompile(`(?i)osteosarcoma`), "1"},
		"rhabdomyosarcoma": {regexp.MustCompile(`(?i)rhabdomyosarcoma`), "1"},
	}
	m.types["other"] = map[string]diagnosis{
		"polyp":                {regexp.MustCompile(`(?i)polyp`), "0"},
		"cyst":                 {regexp.MustCompile(`(?i)cyst`), "0"},
		"papilloma":            {regexp.MustCompile(`(?i)papilloma`), "0"},
		"epulis":               {regexp.MustCompile(`(?i)epuli.*`), "0"},
		"meningioma":           {regexp.MustCompile(`(?i)meningioma`), "0"},
		"hyperplasia":          {regexp.MustCompile(`(?i)(meta|dys|hyper)plas(ia|tic)`), "0"},
		"trichoepithelioma":    {regexp.MustCompile(`(?i)trichoepithelioma`), "0"},
		"lipoma":               {regexp.MustCompile(`(?i)lipoma`), "0"},
		"fibropapilloma ":      {regexp.MustCompile(`(?i)fibropapilloma `), "0"},
		"fibropapillomatosis":  {regexp.MustCompile(`(?i)fibropapillomatosis `), "0"},
		"epithelioma":          {regexp.MustCompile(`(?i)epithelioma`), "0"},
		"leiomyoma":            {regexp.MustCompile(`(?i)leiomyoma`), "0"},
		"hemangioma":           {regexp.MustCompile(`(?i)hemangioma`), "0"},
		"insulinoma":           {regexp.MustCompile(`(?i)insulinoma`), "0"},
		"fibroma":              {regexp.MustCompile(`(?i)fibroma`), "0"},
		"odontoma":             {regexp.MustCompile(`(?i)odontoma`), "0"},
		"osteoma":              {regexp.MustCompile(`(?i)osteoma`), "0"},
		"chromatophroma":       {regexp.MustCompile(`(?i)chromatophroma`), "0"},
		"hamartoma ":           {regexp.MustCompile(`(?i)hamartoma `), "0"},
		"neoplasia":            {regexp.MustCompile(`(?i)neoplasia|neoplasm|tumor`), "-1"},
		"leukemia":             {regexp.MustCompile(`(?i)leukemia`), "1"},
		"chordoma":             {regexp.MustCompile(`(?i)chordoma`), "1"},
		"disseminated":         {regexp.MustCompile(`(?i)disseminated`), "1"},
		"granulosa cell tumor": {regexp.MustCompile(`(?i)Granulosa cell tumor`), "1"},
		"fibrous histiocytoma": {regexp.MustCompile(`(?i)(fibrous )?histiocytoma`), "-1"},
		"melanoma":             {regexp.MustCompile(`(?i)melanoma`), "1"},
		"nephroblastoma":       {regexp.MustCompile(`(?i)(nephroblastoma|(Wilm’s Tumor ))`), "1"},
		"pheochromocytoma":     {regexp.MustCompile(`(?i)pheochromocytoma`), "1"},
		"seminoma":             {regexp.MustCompile(`(?i)seminoma`), "1"},
	}
}
