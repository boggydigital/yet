package rest

import (
	"strings"
)

func writeSharedStyles(sb *strings.Builder) {
	sb.WriteString(
		"body {background: black; color: white;font-family:sans-serif; margin: 2rem;} " +
			"a {display:flex;flex-direction: column; color:white;font-size:1.3rem;font-weight:bold;text-decoration:none;margin-block:2rem;}" +
			"a img {border-radius:0.25rem;width:200px;aspect-ratio:16/9;background:dimgray}" +
			"a.ended {filter:grayscale(1.0)}" +
			"a.highlight {color:gold; margin-block:2rem}" +
			".title {margin-block-start:0.5rem;margin-block-end:0.25rem}" +
			".subtitle {font-size:66%; font-weight:normal}" +
			"a .subtitle {color:dimgray}" +
			"div.subtle {color: dimgray}" +
			"details {margin-block:2rem; content-visibility: auto}" +
			"summary {margin-block-end: 1rem}" +
			"summary::after {content: '\u2026';flex-shrink: 0}" +
			"summary::-webkit-details-marker {display: none}" +
			"summary h1, summary h2 {display: inline; cursor: pointer;color:turquoise}")
}
