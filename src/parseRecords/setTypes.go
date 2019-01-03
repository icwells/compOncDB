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
		"breast":          regexp.MustCompile(`(?i)(breast|mammary)`),
		"colon":           regexp.MustCompile(`(?i)colon|rectum`),
		"duodenum":        regexp.MustCompile(`(?i)duodenum`),
		"fat":             regexp.MustCompile(`(?i)fat|adipose.*`),
		"heart":           regexp.MustCompile(`(?i)heart|cardiac|atrial`),
		"kidney":          regexp.MustCompile(`(?i)kidney.*|ureter|renal`),
		"leukemia":        regexp.MustCompile(`(?i)leukemia`),
		"liver":           regexp.MustCompile(`(?i)hepa.*|liver.*|hep.*|billia.*`),
		"lung":            regexp.MustCompile(`(?i)lung.*|pulm.*|mediasti.*|bronchial|alveol.*`),
		"lymph nodes":     regexp.MustCompile(`(?i)lymph|lymph node`),
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
		"widespread":      regexp.MustCompile(`(?i)widespread|metastatic|(body as a whole)|multiple|disseminated`),
	}
	m.types = map[string]diagnosis{
		{"adenoma": regexp.MustCompile(`(?i)adenoma`), "N"},
		{"polyp": regexp.MustCompile(`(?i)polyp`), "N"},
		{"cyst": regexp.MustCompile(`(?i)cyst`), "N"},
		{"papilloma": regexp.MustCompile(`(?i)papilloma`), "N"},
		{"epulis": regexp.MustCompile(`(?i)epuli.*`), "N"},
		{"meningioma": regexp.MustCompile(`(?i)meningioma`), "N"},
		{"hyperplasia": regexp.MustCompile(`(?i)(meta|dys|hyper)plas(ia|tic)`), "N"},
		{"cystadenoma": regexp.MustCompile(`(?i)cystadenoma`), "N"},
		{"trichoepithelioma": regexp.MustCompile(`(?i)trichoepithelioma`), "N"},
		{"lipoma": regexp.MustCompile(`(?i)lipoma`), "N"},
		{"fibropapilloma ": regexp.MustCompile(`(?i)fibropapilloma `), "N"},
		{"fibropapillomatosis": regexp.MustCompile(`(?i)fibropapilloma `), "N"},
		{"epithelioma": regexp.MustCompile(`(?i)epithelioma`), "N"},
		{"leiomyoma": regexp.MustCompile(`(?i)leiomyoma`), "N"},
		{"hemangioma": regexp.MustCompile(`(?i)hemangioma`), "N"},
		{"insulinoma": regexp.MustCompile(`(?i)insulinoma`), "N"},
		{"fibroma": regexp.MustCompile(`(?i)fibroma`), "N"},
		{"odontoma": regexp.MustCompile(`(?i)odontoma`), "N"},
		{"osteoma": regexp.MustCompile(`(?i)osteoma`), "N"},
		{"chromatophroma": regexp.MustCompile(`(?i)chromatophroma`), "N"},
		{"hamartoma ": regexp.MustCompile(`(?i)hamartoma `), "N"},
		{"neoplasia": regexp.MustCompile(`(?i)neoplasia|neoplasm|tumor`), "NA"},
		{"adenocarcinoma": regexp.MustCompile(`(?i)adenocarcinoma`), "Y"},
		{"carcinoma": regexp.MustCompile(`(?i)carcinoma|TCC`), "Y"},
		{"lymphoma": regexp.MustCompile(`(?i)lymphoma|lymphosarcoma`), "Y"},
		{"leukemia": regexp.MustCompile(`(?i)leukemia`), "Y"},
		{"sarcoma": regexp.MustCompile(`(?i)sarcoma`), "Y"},
		{"carcinomatosis": regexp.MustCompile(`(?i)carinomatosis`), "Y"},
		{"chondrosarcoma": regexp.MustCompile(`(?i)chondrosarcoma`), "Y"},
		{"chordoma": regexp.MustCompile(`(?i)chordoma`), "Y"},
		{"disseminated": regexp.MustCompile(`(?i)disseminated`), "Y"},
		{"fibroadenocarcinoma": regexp.MustCompile(`(?i)fibroadenocarcinoma`), "Y"},
		{"fibrosarcoma": regexp.MustCompile(`(?i)fibrosarcoma`), "Y"},
		{"granulosa cell tumor": regexp.MustCompile(`(?i)Granulosa cell tumor`), "Y"},
		{"hemangiosarcoma": regexp.MustCompile(`(?i)Hemangiosarcoma`), "Y"},
		{"leiomyosarcoma": regexp.MustCompile(`(?i)leiomyosarcoma`), "Y"},
		{"fibrous histiocytoma": regexp.MustCompile(`(?i)(fibrous )?histiocytoma`), "NA"},
		{"melanoma": regexp.MustCompile(`(?i)melanoma`), "Y"},
		{"myxosarcoma": regexp.MustCompile(`(?i)myxosarcoma`), "Y"},
		{"nephroblastoma": regexp.MustCompile(`(?i)(nephroblastoma|(Wilm’s Tumor ))`), "Y"},
		{"osteosarcoma": regexp.MustCompile(`(?i)osteosarcoma`), "Y"},
		{"pheochromocytoma": regexp.MustCompile(`(?i)pheochromocytoma`), "Y"},
		{"rhabdomyosarcoma": regexp.MustCompile(`(?i)rhabdomyosarcoma`), "Y"},
		{"seminoma": regexp.MustCompile(`(?i)seminoma`), "Y"},
	}
}
