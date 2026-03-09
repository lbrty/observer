---
title: Berichte
weight: 7
---

Observer generiert 39 strukturierte Berichtstypen, die den tatsächlichen Berichtspflichten ukrainischer NGOs gegenüber Geldgebern entsprechen. Alle Berichte sind auf ein einzelnes Projekt beschränkt und werden nach einem Datumsbereich gefiltert.

## Einen Bericht ausführen

Alle Berichte sind verfügbar unter:

```http
GET /projects/:project_id/reports/:report_type?start=2024-01-01&end=2024-12-31
```

Ersetzen Sie `:report_type` durch einen der Bezeichner aus den Tabellen unten.

## Zähltypen

Jeder Bericht verwendet eine von drei Zählmethoden:

| Typ         | Bedeutung                                                                   |
| ----------- | --------------------------------------------------------------------------- |
| **Ereignisse** | Zeilen in `support_records` — eine Person mit drei Beratungen zählt als 3 |
| **Personen** | Einzigartige Individuen — dieselbe Person wird einmal gezählt, unabhängig von der Anzahl der Datensätze |
| **Einheiten** | Einzigartige Familieneinheiten (ein Haushalt = eine Einheit)               |

Die Datumsfilterung verwendet `support_records.provided_at` für Beratungen und `people.registered_at` für Registrierungsberichte — nicht `created_at`.

## Berichtsgruppen

### Gruppe 1 — Allgemeine Beratungszählungen

| # | Bericht | Typ |
|---|---------|-----|
| 1 | Gesamtberatungen aller Typen | Ereignisse |
| 2 | Gesamtrechtsberatungen | Ereignisse |
| 3 | Gesamtsozialberatungen | Ereignisse |

### Gruppe 2 — Geschlechteraufteilung

| # | Bericht | Typ |
|---|---------|-----|
| 12 | Im Zeitraum registrierte Männer | Personen |
| 13 | Im Zeitraum registrierte Frauen | Personen |
| 14 | Frauen, die Rechtsberatungen erhielten | Personen |
| 15 | Frauen, die Sozialberatungen erhielten | Personen |
| 16 | Männer, die Rechtsberatungen erhielten | Personen |
| 17 | Männer, die Sozialberatungen erhielten | Personen |

Berichte 14–17 zählen einzigartige Personen, keine Beratungsereignisse.

### Gruppe 3 — Geografischer / IDP-Status

| # | Bericht | Typ |
|---|---------|-----|
| 4 | Insgesamt registrierte Personen | Personen |
| 5–6 | Aus Konfliktgebieten registrierte Personen | Personen |
| 7–10 | Personen aus Konfliktgebieten, die Rechts-/Sozialberatungen erhielten | Personen |
| 11 | Registrierte Nicht-IDPs | Personen |

Der IDP-Status wird automatisch abgeleitet: `origin_place_id → places → states.conflict_zone`. Berichte 5–6 und 7–10 sind nach Konfliktzonenbezeichnung parametrisiert.

### Gruppe 4 — Aufschlüsselung nach Vulnerabilitätskategorie

| # | Bericht | Typ |
|---|---------|-----|
| 18 | Registrierte Personen — nach Vulnerabilitätskategorie | Personen |
| 19 | Personen, die Sozialberatungen erhielten — nach Kategorie | Personen |
| 20 | Personen, die Rechtsberatungen erhielten — nach Kategorie | Personen |

Personen ohne zugewiesene Kategorie erscheinen in einem „unkategorisiert"-Bucket.

### Gruppe 5 — Aktuelle Aufenthaltsregion

| # | Bericht | Typ |
|---|---------|-----|
| 21 | Registrierte Personen — nach aktueller Region | Personen |
| 22 | Personen, die Rechtsberatungen erhielten — nach Region | Personen |
| 23 | Personen, die Sozialberatungen erhielten — nach Region | Personen |

### Gruppe 6 — Aufschlüsselung nach Unterstützungsbereich

| # | Bericht | Typ |
|---|---------|-----|
| 24 | Rechtsberatungsanzahl — nach Bereich | Ereignisse |
| 25 | Personen, die Rechtsberatungen erhielten — nach Bereich | Personen |
| 29 | Sozialberatungsanzahl — nach Bereich | Ereignisse |
| 30 | Personen, die Sozialberatungen erhielten — nach Bereich | Personen |

**Unterstützungsbereiche:**

| Wert | Beschreibung |
|------|--------------|
| `housing_assistance` | Wohnrechte, Zwangsräumung, Sozialwohnungen |
| `document_recovery` | Pässe, Geburtsurkunden, Eigentumsunterlagen |
| `social_benefits` | IDP-Registrierung, Sozialleistungen |
| `property_rights` | In besetzten Gebieten zurückgelassenes Eigentum |
| `employment_rights` | Arbeitsrecht, Entlassung, Arbeitsvermittlung |
| `family_law` | Scheidung, Sorgerecht, Unterhalt |
| `healthcare_access` | Krankenversicherung, Behinderungsdokumentation |
| `education_access` | Schulanmeldung, Bildungsrechte |
| `financial_aid` | Finanzielle Soforthilfe |
| `psychological_support` | Überweisungen an Fachkräfte für psychische Gesundheit, Beratung |
| `other` | Nicht aufgeführte oder bereichsübergreifende Themen |

### Gruppe 7 — Aufschlüsselung nach Büro

| # | Bericht | Typ |
|---|---------|-----|
| 28 | Rechtsberatungsanzahl — nach Büro | Ereignisse |
| 32 | Sozialberatungsanzahl — nach Büro | Ereignisse |
| 33 | Gesamtberatungsanzahl — nach Büro | Ereignisse |

### Gruppe 8 — Aufschlüsselung nach Altersgruppe

| # | Bericht | Typ |
|---|---------|-----|
| 26 | Rechtsberatungsanzahl — nach Altersgruppe | Ereignisse |
| 27 | Personen, die Rechtsberatungen erhielten — nach Altersgruppe | Personen |
| 31a | Sozialberatungsanzahl — nach Altersgruppe | Ereignisse |
| 31b | Personen, die Sozialberatungen erhielten — nach Altersgruppe | Personen |
| 34 | Gesamtberatungsanzahl — nach Altersgruppe | Ereignisse |

**Altersgruppen:**

| Wert | Altersbereich |
|------|---------------|
| `infant` | 0–1 |
| `toddler` | 1–3 |
| `pre_school` | 3–6 |
| `middle_childhood` | 6–12 |
| `young_teen` | 12–14 |
| `teenager` | 14–18 |
| `young_adult` | 18–25 |
| `early_adult` | 25–35 |
| `middle_aged_adult` | 35–55 |
| `old_adult` | 55+ |

Wenn `birth_date` gesetzt ist und `age_group` null ist, berechnet die Anwendung den Bucket automatisch.

### Gruppe 9 — Tag-Suche

| # | Bericht | Typ |
|---|---------|-----|
| 35 | Unterstützungsdatensätze für Personen mit bestimmten Tags | Ereignisse |
| 36 | Mit bestimmten Tags registrierte Personen | Personen |

Übergeben Sie eine oder mehrere Tag-IDs als Parameter. Nützlich für Ad-hoc-Geberanfragen.

### Gruppe 10 — Familieneinheiten

| # | Bericht | Typ |
|---|---------|-----|
| 37 | Personen und Familienmitglieder, die Rechtsberatungen erhielten | Personen + Einheiten |
| 38 | Personen und Familienmitglieder, die Sozialberatungen erhielten | Personen + Einheiten |
| 39 | Im Zeitraum registrierte Personen und Familienmitglieder | Personen + Einheiten |

Eine Familieneinheit ist ein Haushaltsdatensatz. Die Zählungen werden sowohl als Gesamtzahl der Personen als auch als einzigartige Haushaltseinheiten zurückgegeben.

## Tierberichte

Tierbezogene Berichte sind separat verfügbar unter:

```http
GET /projects/:project_id/pet-reports/:report_type
```

Tierberichte umfassen die Aufschlüsselung nach Tierart, Impfstatus, Sterilisationsstatus und das Verhältnis von Tieren zu Personen.
