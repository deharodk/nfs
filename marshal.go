package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type Doc struct {
	XMLName     xml.Name        `xml:"Comprobante"`
	Tipo        string          `xml:"tipoDeComprobante,attr"`
	Version     string          `xml:"version,attr"`
	Emisor      CFDIEmisor      `xml:"Emisor"`
	Receptor    CFDIReceptor    `xml:"Receptor"`
	Conceptos   []CFDIConcepto  `xml:"Conceptos>Concepto"`
	Complemento CFDIComplemento `xml:"Complemento"`
}

type CFDIEmisor struct {
	XMLName xml.Name `xml:"Emisor"`
	RFC     string   `xml:"rfc,attr"`
}

type CFDIReceptor struct {
	XMLName xml.Name `xml:"Receptor"`
	RFC     string   `xml:"rfc,attr"`
}

type CFDIConcepto struct {
	XMLName     xml.Name `xml:"Concepto"`
	Descripcion string   `xml:"descripcion,attr"`
}

type CFDIComplemento struct {
	XMLName             xml.Name               `xml:"Complemento"`
	TimbreFiscalDigital TFDTimbreFiscalDigital `xml:"TimbreFiscalDigital"`
}

type TFDTimbreFiscalDigital struct {
	XMLName           xml.Name `xml:"TimbreFiscalDigital"`
	NumeroCertificado string   `xml:"noCertificadoSAT,attr"`
	FechaTimbrado     string   `xml:"FechaTimbrado,attr"`
	UUID              string   `xml:"UUID,attr"`
}

func (t TFDTimbreFiscalDigital) String() string {
	return fmt.Sprintf("%s\t%s", t.NumeroCertificado, t.FechaTimbrado)
}

// func (d Doc) String() string {
// 	return fmt.Sprintf("%s\t%s\t%s\t%s\t%s",
// 		d.Emisor.RFC,
// 		d.Receptor.RFC,
// 		d.Complemento.TimbreFiscalDigital.NumeroCertificado,
// 		d.Complemento.TimbreFiscalDigital.FechaTimbrado,
// 		d.Version)
// }

func (c CFDIConcepto) ContainsKeyword() bool {
	desc := strings.ToLower(c.Descripcion)
	return strings.Contains(desc, "magna") ||
		strings.Contains(desc, "premium") ||
		strings.Contains(desc, "diesel")
}

func (d Doc) ContainsGasKeyword() bool {
	for _, concept := range d.Conceptos {
		if concept.ContainsKeyword() {
			return true
		}
	}
	return false
}

func parseXml(doc []byte) Doc {
	var query Doc
	xml.Unmarshal(doc, &query)
	return query
}

func split(path string) [][]byte {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer file.Close()

	rawContent, _ := ioutil.ReadAll(file)
	re := regexp.MustCompile(`(?s)<cfdi:Comprobante\b*[^>]*>(.*?)</cfdi:Comprobante>`)
	matches := re.FindAll(rawContent, -1)
	return matches
}

func splitFile(path string) {
	matches := split(path)
	for _, doc := range matches {
		cfdi := parseXml(doc)
		// fmt.Println(len(matches) != 1)
		fmt.Printf("'%s'\t%s\t%d\n",
			cfdi.Complemento.TimbreFiscalDigital.UUID,
			cfdi.Emisor.RFC,
			len(matches))
	}
}
