package handlers

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl-lang/decoder"
	"github.com/hashicorp/hcl-lang/lang"
	lsctx "github.com/hashicorp/terraform-ls/internal/context"
	ilsp "github.com/hashicorp/terraform-ls/internal/lsp"
	"github.com/sourcegraph/go-lsp"
	"github.com/zclconf/go-cty/cty"
)

type DocumentSymbol struct {
	Name           string           `json:"name"`
	Detail         string           `json:"detail,omitempty"`
	Kind           lsp.SymbolKind   `json:"kind"`
	Deprecated     bool             `json:"deprecated,omitempty"`
	Range          lsp.Range        `json:"range"`
	SelectionRange lsp.Range        `json:"selectionRange"`
	Children       []DocumentSymbol `json:"children,omitempty"`
}

func (h *logHandler) TextDocumentSymbol(ctx context.Context, params lsp.DocumentSymbolParams) ([]DocumentSymbol, error) {
	var symbols []DocumentSymbol

	fs, err := lsctx.DocumentStorage(ctx)
	if err != nil {
		return symbols, err
	}

	df, err := lsctx.DecoderFinder(ctx)
	if err != nil {
		return symbols, err
	}

	file, err := fs.GetDocument(ilsp.FileHandlerFromDocumentURI(params.TextDocument.URI))
	if err != nil {
		return symbols, err
	}

	// TODO: block until it's available <-df.ParserLoadingDone()
	// requires https://github.com/hashicorp/terraform-ls/issues/8
	// textDocument/documentSymbol fires early alongside textDocument/didOpen
	// the protocol does not retry the request, so it's best to give the parser
	// a moment
	if err := Waiter(func() (bool, error) {
		return df.IsCoreSchemaLoaded(file.Dir())
	}).Waitf("core schema is not available yet for %s", file.Dir()); err != nil {
		return symbols, err
	}

	d, err := df.DecoderForDir(file.Dir())
	if err != nil {
		return symbols, fmt.Errorf("finding compatible decoder failed: %w", err)
	}

	sbs, err := d.Symbols()
	if err != nil {
		return symbols, err
	}

	return convertSymbols(sbs), nil
}

func convertSymbols(sbs []decoder.Symbol) []DocumentSymbol {
	symbols := make([]DocumentSymbol, len(sbs))
	for i, s := range sbs {
		var kind lsp.SymbolKind
		switch s.Kind() {
		case lang.BlockSymbolKind:
			kind = lsp.SKClass
		case lang.AttributeSymbolKind:
			kind = attributeSymbolKind(s)
		}

		symbols[i] = DocumentSymbol{
			Name:           s.Name(),
			Kind:           kind,
			Range:          ilsp.HCLRangeToLSP(s.Range()),
			SelectionRange: ilsp.HCLRangeToLSP(s.Range()),
			Children:       convertSymbols(s.NestedSymbols()),
		}
	}
	return symbols
}

func attributeSymbolKind(s decoder.Symbol) lsp.SymbolKind {
	as, ok := s.(*decoder.AttributeSymbol)
	if !ok {
		return lsp.SKProperty
	}

	switch as.Type {
	case cty.String:
		return lsp.SKString
	case cty.Number:
		return lsp.SKNumber
	case cty.Bool:
		return lsp.SKBoolean
	}

	if as.Type.IsListType() || as.Type.IsSetType() {
		return lsp.SKArray
	}

	if as.Type.IsObjectType() {
		return lsp.SKObject
	}

	return lsp.SKField
}
