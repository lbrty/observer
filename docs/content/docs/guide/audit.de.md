---
title: Prüfprotokoll
weight: 9
---

Observer zeichnet einen Prüfprotokoll-Eintrag für jeden Erstellungs-, Aktualisierungs- und Löschvorgang an Falldaten auf. Prüfprotokolle sind nur anfügbar — sie können nicht bearbeitet oder gelöscht werden.

## Was wird protokolliert

| Kategorie          | Protokollierte Vorgänge                                   |
| ------------------ | --------------------------------------------------------- |
| Personen           | Erstellen, aktualisieren, löschen von Personendatensätzen |
| Unterstützungsakte | Erstellen, aktualisieren, löschen von Beratungen          |
| Migrationsakte     | Erstellen, aktualisieren, löschen von Bewegungsdatensätzen |
| Haushalte          | Erstellen, aktualisieren, löschen von Haushalten und Mitgliedern |
| Notizen            | Erstellen, aktualisieren, löschen von Fallnotizen         |
| Dokumente          | Hochladen, Metadaten aktualisieren, Dokumente löschen     |
| Haustiere          | Erstellen, aktualisieren, löschen von Haustierdatensätzen |
| Berechtigungen     | Vergeben, aktualisieren, entziehen von Projektberechtigungen |

Authentifizierungsereignisse (Anmeldung, Abmeldung, Token-Aktualisierung) befinden sich nicht im Projekt-Prüfprotokoll — sie erscheinen in den Serverprotokollen.

## Prüfprotokoll anzeigen

Nur Projektmanager und Eigentümer können auf das Prüfprotokoll zugreifen.

```http
GET /projects/:project_id/audit?page=1&per_page=50
```

| Parameter  | Beschreibung                                    |
| ---------- | ----------------------------------------------- |
| `page`     | Seitennummer (Standard 1)                       |
| `per_page` | Ergebnisse pro Seite (Standard 50)              |
| `actor_id` | Nach dem Benutzer filtern, der die Änderung vorgenommen hat |
| `start`    | Nach Datum filtern (JJJJ-MM-TT)                |
| `end`      | Nach Datum filtern (JJJJ-MM-TT)                |

### Antwortformat

Jeder Prüfeintrag enthält:

```json
{
  "id": "01J...",
  "project_id": "01J...",
  "actor_id": "01J...",
  "actor_ip": "192.168.1.1",
  "action": "create",
  "entity_type": "person",
  "entity_id": "01J...",
  "created_at": "2024-06-15T10:30:00Z"
}
```

| Feld          | Beschreibung                                                |
| ------------- | ----------------------------------------------------------- |
| `actor_id`    | Benutzer, der die Aktion durchgeführt hat                   |
| `actor_ip`    | IP-Adresse der Anfrage                                      |
| `action`      | `create`, `update` oder `delete`                            |
| `entity_type` | Datensatztyp (`person`, `support_record`, `note` usw.)     |
| `entity_id`   | ULID des betroffenen Datensatzes                            |

## Prüfung im Web-Interface

1. Ein Projekt öffnen
2. In der Navigation auf **Prüfprotokoll** klicken
3. Nach Datumsbereich oder Benutzer filtern

{{% callout type="note" %}}
Der Zugriff auf das Prüfprotokoll erfordert die Projektrolle `manager` oder `owner`. Berater und Betrachter können das Prüfprotokoll nicht einsehen.
{{% /callout %}}
