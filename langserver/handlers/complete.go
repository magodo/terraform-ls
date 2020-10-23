package handlers

import (
	"context"
	"fmt"

	lsctx "github.com/hashicorp/terraform-ls/internal/context"
	ilsp "github.com/hashicorp/terraform-ls/internal/lsp"
	lsp "github.com/sourcegraph/go-lsp"
)

func (h *logHandler) TextDocumentComplete(ctx context.Context, params lsp.CompletionParams) (lsp.CompletionList, error) {
	var list lsp.CompletionList

	fs, err := lsctx.DocumentStorage(ctx)
	if err != nil {
		return list, err
	}

	cc, err := lsctx.ClientCapabilities(ctx)
	if err != nil {
		return list, err
	}

	df, err := lsctx.DecoderFinder(ctx)
	if err != nil {
		return list, err
	}

	file, err := fs.GetDocument(ilsp.FileHandlerFromDocumentURI(params.TextDocument.URI))
	if err != nil {
		return list, err
	}

	isReady, err := df.IsDecoderReady(file.Dir())
	if err != nil {
		return list, err
	}
	if !isReady {
		return list, fmt.Errorf("decoder is not available yet for %s", file.Dir())
	}

	isCoreSchemaLoaded, err := df.IsCoreSchemaLoaded(file.Dir())
	if err != nil {
		return list, err
	}
	if !isCoreSchemaLoaded {
		// TODO: block until it's available <-df.CoreSchemaLoadingDone()
		// requires https://github.com/hashicorp/terraform-ls/issues/8
		return list, fmt.Errorf("core schema is not available yet for %s", file.Dir())
	}

	d, err := df.DecoderForDir(file.Dir())
	if err != nil {
		return list, fmt.Errorf("finding compatible decoder failed: %w", err)
	}

	fPos, err := ilsp.FilePositionFromDocumentPosition(params.TextDocumentPositionParams, file)
	if err != nil {
		return list, err
	}

	h.logger.Printf("Looking for candidates at %q -> %#v", file.Filename(), fPos.Position())
	candidates, diags := d.CandidatesAtPos(file.Filename(), fPos.Position())
	if diags.HasErrors() {
		err = diags
	}

	return ilsp.CompletionList(candidates, cc.TextDocument), err
}
