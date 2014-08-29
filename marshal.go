package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Doc struct {
	XMLName     xml.Name        `xml:"Comprobante"`
	Tipo        string          `xml:"tipoDeComprobante,attr"`
	Version     string          `xml:"version,attr"`
	Serie       string          `xml:"serie,attr"`
	Folio       string          `xml:"folio,attr"`
	Fecha       string          `xml:"fecha,attr"`
	Moneda      string          `xml:"Moneda,attr"`
	TipoCambio  string          `xml:"TipoCambio,attr"`
	Total       string          `xml:"total,attr"`
	SubTotal    string          `xml:"subTotal,attr"`
	Emisor      CFDIEmisor      `xml:"Emisor"`
	Receptor    CFDIReceptor    `xml:"Receptor"`
	Conceptos   []CFDIConcepto  `xml:"Conceptos>Concepto"`
	Impuestos   CFDIImpuestos   `xml:"Impuestos"`
	Complemento CFDIComplemento `xml:"Complemento"`
	Addenda     CFDIAddenda     `xml:"Addenda"`
}

type CFDIImpuestos struct {
	XMLName   xml.Name      `xml:"Impuestos"`
	Total     string        `xml:"totalImpuestosTrasladados,attr"`
	Traslados CFDITraslados `xml:"Traslados"`
}

type CFDITraslados struct {
	XMLName  xml.Name     `xml:"Traslados"`
	Traslado CFDITraslado `xml:"Traslado"`
}

type CFDITraslado struct {
	XMLName xml.Name `xml:"Traslado"`
	Importe string   `xml:"importe,attr"`
}

type CFDIAddenda struct {
	XMLName            xml.Name               `xml:"Addenda"`
	AddendaBuzonFiscal AddendaBuzonFiscalNode `xml:"AddendaBuzonFiscal"`
}

type AddendaBuzonFiscalNode struct {
	XMLName xml.Name `xml:"AddendaBuzonFiscal"`
	CFD     CFDNode  `xml:"CFD"`
}

type CFDNode struct {
	XMLName xml.Name `xml:"CFD"`
	RefID   string   `xml:"refID,attr"`
}

type CFDIEmisor struct {
	XMLName xml.Name `xml:"Emisor"`
	RFC     string   `xml:"rfc,attr"`
}

type CFDIReceptor struct {
	XMLName xml.Name `xml:"Receptor"`
	RFC     string   `xml:"rfc,attr"`
	Nombre  string   `xml:"nombre,attr"`
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

func parseXml(doc []byte) Doc {
	var query Doc
	xml.Unmarshal(doc, &query)
	return query
}

func EncodeAsRow(path string) string {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer file.Close()
	rawContent, _ := ioutil.ReadAll(file)
	cfdi := parseXml(rawContent)

	row := []string{
		cfdi.Addenda.AddendaBuzonFiscal.CFD.RefID,
		cfdi.Receptor.RFC,
		cfdi.Receptor.Nombre,
		cfdi.Serie,
		cfdi.Folio,
		cfdi.Fecha,
		cfdi.Moneda,
		cfdi.TipoCambio,
		cfdi.Total,
		cfdi.Impuestos.Traslados.Traslado.Importe,
		cfdi.SubTotal,
		cfdi.Tipo,
		cfdi.Complemento.TimbreFiscalDigital.FechaTimbrado,
		cfdi.Complemento.TimbreFiscalDigital.UUID}
	return strings.Join(row, "\t")
}
