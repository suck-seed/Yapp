package services

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"net/url"
// 	"os"
// 	"strconv"
// 	"strings"
// 	"time"
// )

// type ModerationAnalysis struct {
// 	ShouldBlock     bool
// 	ShouldIncrement bool

// 	Severity int
// 	Reason   string

// 	MatchedLabels []string
// 	MatchedTerms  []string
// 	Scores        map[string]float64
// }

// type ITextModerationDetector interface {
// 	Analyze(ctx context.Context, content string) (*ModerationAnalysis, error)
// }

// type SightengineDetector struct {
// 	apiUser   string
// 	apiSecret string
// 	lang      string
// 	client    *http.Client
// 	enabled   bool
// }

// func NewSightengineDetectorFromEnv() ITextModerationDetector {
// 	enabled := strings.EqualFold(os.Getenv("SIGHTENGINE_ENABLED"), "true")
// 	apiUser := os.Getenv("SIGHTENGINE_API_USER")
// 	apiSecret := os.Getenv("SIGHTENGINE_API_SECRET")

// 	if !enabled || apiUser == "" || apiSecret == "" {
// 		log.Println("Sightengine moderation disabled or missing credentials")
// 		return &NoopTextModerationDetector{}
// 	}

// 	lang := os.Getenv("SIGHTENGINE_TEXT_LANG")
// 	if strings.TrimSpace(lang) == "" {
// 		lang = "en"
// 	}

// 	return &SightengineDetector{
// 		apiUser:   apiUser,
// 		apiSecret: apiSecret,
// 		lang:      lang,
// 		enabled:   true,
// 		client: &http.Client{
// 			Timeout: 1500 * time.Millisecond,
// 		},
// 	}
// }

// type NoopTextModerationDetector struct{}

// func (d *NoopTextModerationDetector) Analyze(ctx context.Context, content string) (*ModerationAnalysis, error) {
// 	return &ModerationAnalysis{
// 		ShouldBlock:     false,
// 		ShouldIncrement: false,
// 		Severity:        0,
// 		Reason:          "moderation_disabled",
// 		Scores:          map[string]float64{},
// 	}, nil
// }

// type sightengineTextResponse struct {
// 	Status string `json:"status"`

// 	Error *struct {
// 		Code    any    `json:"code"`
// 		Message string `json:"message"`
// 	} `json:"error,omitempty"`

// 	Profanity struct {
// 		Matches []sightengineProfanityMatch `json:"matches"`
// 	} `json:"profanity"`

// 	ModerationClasses map[string]any `json:"moderation_classes"`
// }

// type sightengineProfanityMatch struct {
// 	Type      string `json:"type"`
// 	Match     string `json:"match"`
// 	Intensity string `json:"intensity"`
// 	Start     int    `json:"start"`
// 	End       int    `json:"end"`
// }

// func (d *SightengineDetector) Analyze(ctx context.Context, content string) (*ModerationAnalysis, error) {
// 	content = strings.TrimSpace(content)
// 	if content == "" {
// 		return &ModerationAnalysis{
// 			Scores: map[string]float64{},
// 		}, nil
// 	}

// 	form := url.Values{}
// 	form.Set("text", content)
// 	form.Set("lang", d.lang)

// 	// rules = profanity/obfuscation
// 	// ml = context-aware classes
// 	form.Set("mode", "rules,ml")
// 	form.Set("models", "general")

// 	form.Set("api_user", d.apiUser)
// 	form.Set("api_secret", d.apiSecret)

// 	req, err := http.NewRequestWithContext(
// 		ctx,
// 		http.MethodPost,
// 		"https://api.sightengine.com/1.0/text/check.json",
// 		strings.NewReader(form.Encode()),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

// 	res, err := d.client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer res.Body.Close()

// 	if res.StatusCode < 200 || res.StatusCode >= 300 {
// 		return nil, fmt.Errorf("sightengine returned http status %d", res.StatusCode)
// 	}

// 	var out sightengineTextResponse
// 	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
// 		return nil, err
// 	}

// 	if out.Status != "success" {
// 		if out.Error != nil {
// 			return nil, fmt.Errorf("sightengine error: %s", out.Error.Message)
// 		}
// 		return nil, fmt.Errorf("sightengine error: status=%s", out.Status)
// 	}

// 	return decideSightengineModeration(out), nil
// }

// func decideSightengineModeration(out sightengineTextResponse) *ModerationAnalysis {
// 	analysis := &ModerationAnalysis{
// 		ShouldBlock:     false,
// 		ShouldIncrement: false,
// 		Severity:        0,
// 		Reason:          "clean",
// 		MatchedLabels:   []string{},
// 		MatchedTerms:    []string{},
// 		Scores:          map[string]float64{},
// 	}

// 	sexual := score(out.ModerationClasses, "sexual")
// 	discriminatory := score(out.ModerationClasses, "discriminatory")
// 	insulting := score(out.ModerationClasses, "insulting")
// 	violent := score(out.ModerationClasses, "violent")
// 	toxic := score(out.ModerationClasses, "toxic")

// 	analysis.Scores["sexual"] = sexual
// 	analysis.Scores["discriminatory"] = discriminatory
// 	analysis.Scores["insulting"] = insulting
// 	analysis.Scores["violent"] = violent
// 	analysis.Scores["toxic"] = toxic

// 	// ML context-aware blocking.
// 	if violent >= 0.65 {
// 		markBlocked(analysis, 5, "violent_text_detected", "violent")
// 	}

// 	if discriminatory >= 0.70 {
// 		markBlocked(analysis, 5, "discriminatory_text_detected", "discriminatory")
// 	}

// 	if insulting >= 0.75 {
// 		markBlocked(analysis, 4, "insulting_text_detected", "insulting")
// 	}

// 	if toxic >= 0.85 {
// 		markBlocked(analysis, 4, "toxic_text_detected", "toxic")
// 	}

// 	if sexual >= 0.90 {
// 		markBlocked(analysis, 3, "sexual_text_detected", "sexual")
// 	}

// 	mediumOrHighRuleMatches := 0

// 	for _, match := range out.Profanity.Matches {
// 		label := "profanity:" + match.Type + ":" + match.Intensity

// 		analysis.MatchedLabels = append(analysis.MatchedLabels, label)
// 		analysis.MatchedTerms = append(analysis.MatchedTerms, match.Match)

// 		intensity := strings.ToLower(match.Intensity)
// 		matchType := strings.ToLower(match.Type)

// 		switch intensity {
// 		case "high":
// 			markBlocked(analysis, 4, "high_intensity_profanity", label)

// 		case "medium":
// 			mediumOrHighRuleMatches++

// 			// Medium insult/discriminatory/sexual is usually worth blocking.
// 			if matchType == "insult" || matchType == "discriminatory" || matchType == "sexual" {
// 				markBlocked(analysis, 3, "targeted_or_sensitive_profanity", label)
// 			}

// 		case "low":
// 			// Low intensity = mild curse. Let it pass unless ML already blocked it.
// 			if analysis.Reason == "clean" {
// 				analysis.Reason = "low_intensity_profanity_allowed"
// 			}
// 		}
// 	}

// 	// Repeated medium profanity becomes block-worthy.
// 	if mediumOrHighRuleMatches >= 2 {
// 		markBlocked(analysis, 3, "repeated_medium_profanity", "repeated_profanity")
// 	}

// 	analysis.ShouldIncrement = analysis.ShouldBlock
// 	return analysis
// }

// func markBlocked(a *ModerationAnalysis, severity int, reason string, label string) {
// 	a.ShouldBlock = true
// 	a.ShouldIncrement = true

// 	if severity > a.Severity {
// 		a.Severity = severity
// 		a.Reason = reason
// 	}

// 	if label != "" {
// 		a.MatchedLabels = append(a.MatchedLabels, label)
// 	}
// }

// func score(classes map[string]any, key string) float64 {
// 	if classes == nil {
// 		return 0
// 	}

// 	raw, ok := classes[key]
// 	if !ok {
// 		return 0
// 	}

// 	switch v := raw.(type) {
// 	case float64:
// 		return v
// 	case float32:
// 		return float64(v)
// 	case int:
// 		return float64(v)
// 	case string:
// 		parsed, _ := strconv.ParseFloat(v, 64)
// 		return parsed
// 	default:
// 		return 0
// 	}
// }
