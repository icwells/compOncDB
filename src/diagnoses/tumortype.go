// Defines tumorype struct and methods

package diagnoses

/*type diagnosis struct {
	expression *regexp.Regexp
	malignant  string
}

func (m *Matcher) setTypes() {
	// Sets type and location maps
	m.location = map[string]*regexp.Regexp{
		"abdomen":         regexp.MustCompile(`(?i)abdom(e|i).*|omentum|diaphragm`),
		"adrenal":         regexp.MustCompile(`(?i)adrenal|pheochromocytoma`),
		"bile duct":       regexp.MustCompile(`(?i)bile.*|biliary`),
		"bladder":         regexp.MustCompile(`(?i)bladder`),
		"bone":            regexp.MustCompile(`(?i)sacrum|bone.*`),
		"brain":           regexp.MustCompile(`(?i)brain`),
		"breast":          regexp.MustCompile(`(?i)breast|mammary`),
		"colon":           regexp.MustCompile(`(?i)colon|rectum`),
		"duodenum":        regexp.MustCompile(`(?i)duodenum`),
		"fat":             regexp.MustCompile(`(?i)fat|adipose.*`),
		"heart":           regexp.MustCompile(`(?i)heart|cardiac|atrial`),
		"kidney":          regexp.MustCompile(`(?i)kidney.*|ureter|renal`),
		"leukemia":        regexp.MustCompile(`(?i)leukemia`),
		"liver":           regexp.MustCompile(`(?i)hepa.*|liver.*|hep.*|billia.*`),
		"lung":            regexp.MustCompile(`(?i)lung.*|pulm.*|mediasti.*|bronchial|alveol.*`),
		"lymph nodes":     regexp.MustCompile(`(?i)lymph( node)?`),
		"muscle":          regexp.MustCompile(`(?i)muscle|.*structure.*`),
		"nerve":           regexp.MustCompile(`(?i)nerve.*`),
		"other":           regexp.MustCompile(`(?i)gland|basal.*|islet|multifocal|neck|nasal|neuroendo.*`),
		"oral":            regexp.MustCompile(`(?i)oral|tongue|mouth|lip|palate|pharyn.*|laryn.*|gingival`),
		"ovary":           regexp.MustCompile(`(?i)ovar.*`),
		"pancreas":        regexp.MustCompile(`(?i)pancreas.*|islet`),
		"seminal vesicle": regexp.MustCompile(`(?i)seminal vesicle`),
		"skin":            regexp.MustCompile(`(?i)skin|eyelid|(sub)?cutan.*|derm.*`),
		"small intestine": regexp.MustCompile(`(?i)(small )?intestin(e|al)`),
		"spinal cord":     regexp.MustCompile(`(?i)spinal|spine`),
		"spleen":          regexp.MustCompile(`(?i)spleen`),
		"testis":          regexp.MustCompile(`(?i)test(i|e).*`),
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
		"chondrosarcoma":   {regexp.MustCompile(`(?i)chondrosarcoma`), "1"},
		"fibrosarcoma":     {regexp.MustCompile(`(?i)fibrosarcoma`), "1"},
		"hemangiosarcoma":  {regexp.MustCompile(`(?i)Hemangiosarcoma`), "1"},
		"leiomyosarcoma":   {regexp.MustCompile(`(?i)leiomyosarcoma`), "1"},
		"lymphoma":         {regexp.MustCompile(`(?i)lymphoma|lymphosarcoma`), "1"},
		"myxosarcoma":      {regexp.MustCompile(`(?i)myxosarcoma`), "1"},
		"osteosarcoma":     {regexp.MustCompile(`(?i)osteosarcoma`), "1"},
		"rhabdomyosarcoma": {regexp.MustCompile(`(?i)rhabdomyosarcoma`), "1"},
		"sarcoma":          {regexp.MustCompile(`(?i)sarcoma`), "1"},
	}
	m.types["fibropapilloma"] = map[string]diagnosis{
		"fibropapilloma":      {regexp.MustCompile(`(?i)fibropapilloma`), "0"},
		"fibropapillomatosis": {regexp.MustCompile(`(?i)fibropapillomatosis`), "0"},
	}
	m.types["other"] = map[string]diagnosis{
		"chordoma":             {regexp.MustCompile(`(?i)chordoma`), "1"},
		"chromatophroma":       {regexp.MustCompile(`(?i)chromatophroma`), "0"},
		"cyst":                 {regexp.MustCompile(`(?i)cyst`), "0"},
		"disseminated":         {regexp.MustCompile(`(?i)disseminated`), "1"},
		"epithelioma":          {regexp.MustCompile(`(?i)epithelioma`), "0"},
		"epulis":               {regexp.MustCompile(`(?i)epuli.*`), "0"},
		"fibroma":              {regexp.MustCompile(`(?i)fibroma`), "0"},
		"fibrous histiocytoma": {regexp.MustCompile(`(?i)(fibrous )?histiocytoma`), "-1"},
		"granulosa cell tumor": {regexp.MustCompile(`(?i)granulosa cell tumor`), "1"},
		"hamartoma ":           {regexp.MustCompile(`(?i)hamartoma`), "0"},
		"hemangioma":           {regexp.MustCompile(`(?i)hemangioma`), "0"},
		"hyperplasia":          {regexp.MustCompile(`(?i)(meta|dys|hyper)plas(ia|tic)`), "0"},
		"insulinoma":           {regexp.MustCompile(`(?i)insulinoma`), "0"},
		"leiomyoma":            {regexp.MustCompile(`(?i)leiomyoma`), "0"},
		"leukemia":             {regexp.MustCompile(`(?i)leukemia`), "1"},
		"lipoma":               {regexp.MustCompile(`(?i)lipoma`), "0"},
		"melanoma":             {regexp.MustCompile(`(?i)melanoma`), "1"},
		"meningioma":           {regexp.MustCompile(`(?i)meningioma`), "0"},
		"neoplasia":            {regexp.MustCompile(`(?i)neoplasia|neoplasm|tumor`), "-1"},
		"nephroblastoma":       {regexp.MustCompile(`(?i)(nephroblastoma|(Wilm’s Tumor ))`), "1"},
		"odontoma":             {regexp.MustCompile(`(?i)odontoma`), "0"},
		"osteoma":              {regexp.MustCompile(`(?i)osteoma`), "0"},
		"papilloma":            {regexp.MustCompile(`(?i)papilloma`), "0"},
		"pheochromocytoma":     {regexp.MustCompile(`(?i)pheochromocytoma`), "1"},
		"polyp":                {regexp.MustCompile(`(?i)polyp`), "0"},
		"seminoma":             {regexp.MustCompile(`(?i)seminoma`), "1"},
		"trichoepithelioma":    {regexp.MustCompile(`(?i)trichoepithelioma`), "0"},
	}
}*/

