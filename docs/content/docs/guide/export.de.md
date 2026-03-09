---
title: Export
weight: 8
---

Observer kann Falldaten und Berichte als CSV-Dateien exportieren. Der Export erfordert mindestens die Projektrolle `consultant`.

## Falldaten exportieren

Alle Personen und ihre zugehörigen Datensätze für ein Projekt exportieren:

```http
GET /projects/:project_id/export/people?format=csv
```

Optionale Abfrageparameter:

| Parameter     | Beschreibung                                     |
| ------------- | ------------------------------------------------ |
| `start`       | Nach Registrierungsdatum filtern (JJJJ-MM-TT)   |
| `end`         | Nach Registrierungsdatum filtern (JJJJ-MM-TT)   |
| `category_id` | Nach Schwachstellenkategorie filtern             |
| `tag_id`      | Nach Tag filtern                                 |

Die Antwort wird gestreamt — große Datensätze werden schrittweise gesendet, damit die Verbindung nicht abbricht.

## Unterstützungsdatensätze exportieren

Beratungsdatensätze für ein Projekt exportieren:

```http
GET /projects/:project_id/export/support-records?format=csv&start=2024-01-01&end=2024-12-31
```

| Parameter | Beschreibung                                       |
| --------- | -------------------------------------------------- |
| `start`   | Nach `provided_at`-Datum filtern (JJJJ-MM-TT)     |
| `end`     | Nach `provided_at`-Datum filtern (JJJJ-MM-TT)     |
| `type`    | `legal` oder `social`                              |
| `sphere`  | Unterstützungsbereich (z.B. `housing_assistance`)  |

## Migrationsdatensätze exportieren

Bewegungs-/Vertreibungsgeschichte exportieren:

```http
GET /projects/:project_id/export/migration-records?format=csv
```

## CSV-Format

Alle Exporte verwenden UTF-8 CSV mit einer Kopfzeile. Die Felder folgen derselben Struktur wie die API-Antworten. Sensible Felder (`contact`, `personal`, `documents`) werden basierend auf Ihren Projektberechtigungs-Flags ein- oder ausgeschlossen — dieselben Regeln, die für API-Lesevorgänge gelten, gelten auch für Exporte.

## Im Web-Interface herunterladen

1. Ein Projekt öffnen
2. In der Navigation auf **Export** klicken
3. Den Datensatztyp und den Datumsbereich auswählen
4. Auf **CSV herunterladen** klicken

Die Datei wird direkt in Ihren Browser heruntergeladen.

{{% callout type="note" %}}
Exporte sind auf Ihr Projekt beschränkt. Sie können keine Daten aus einem Projekt exportieren, dem Sie nicht zugewiesen sind, und Ihre Sensitivitäts-Flags werden angewendet — wenn `can_view_personal` deaktiviert ist, enthält die exportierte CSV keine nationalen Kennziffern oder Geburtsdaten.
{{% /callout %}}
