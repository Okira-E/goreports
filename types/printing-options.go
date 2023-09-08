package types

type PrintingOptions struct {
	PaperSize    string             `json:"paperSize"`
	Landscape    bool               `json:"landscape"`
	MarginTop    int                `json:"marginTop"`
	MarginRight  int                `json:"marginRight"`
	MarginBottom int                `json:"marginBottom"`
	MarginLeft   int                `json:"marginLeft"`
	PageNumbers  PageNumbersOptions `json:"pageNumbers"`
}

type PageNumbersOptions struct {
	Enabled  bool   `json:"enabled"`
	Position string `json:"position"`
}
