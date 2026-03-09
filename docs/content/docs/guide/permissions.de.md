---
title: Berechtigungen
weight: 6
---

Observer verwendet ein zweistufiges Berechtigungsmodell. Jeder Benutzer hat eine **Plattformrolle**, die den Zugriff auf Verwaltungsfunktionen steuert. Innerhalb jedes Projekts haben Benutzer zusätzlich eine **Projektrolle**, die steuert, was sie mit Falldaten tun können. Drei Vertraulichkeits-Flags schränken zusätzlich ein, welche Felder in Antworten erscheinen.

## Plattformrollen

Plattformrollen werden bei der Erstellung des Benutzerkontos festgelegt und können nur von einem Administrator geändert werden.

| Rolle        | Was sie tun können                                                         |
| ------------ | -------------------------------------------------------------------------- |
| `admin`      | Vollzugriff auf alles — Benutzer, Projekte, Referenzdaten, alle Fälle      |
| `staff`      | Referenzdaten und Projekte verwalten; kann keine anderen Benutzer verwalten |
| `consultant` | In zugewiesenen Projekten arbeiten; kein Zugriff auf das Admin-Panel       |
| `guest`      | Lesezugriff auf zugewiesene Projekte                                       |

Plattformadministratoren agieren automatisch als **Projekteigentümer** in allen Projekten — sie umgehen projektbezogene Berechtigungsprüfungen vollständig.

## Projektrollen

Projektrollen gelten innerhalb eines bestimmten Projekts. Ein Benutzer muss einem Projekt explizit zugewiesen werden, um darauf zugreifen zu können.

| Rolle        | Rang | Was sie tun können                                      |
| ------------ | ---- | ------------------------------------------------------- |
| `owner`      | 4    | Alles, einschließlich Löschen des Projekts              |
| `manager`    | 3    | Alle Falldatenoperationen + Teammitglieder verwalten    |
| `consultant` | 2    | Fälle erstellen, lesen, aktualisieren und Daten exportieren |
| `viewer`     | 1    | Lesezugriff auf Falldaten                               |

### Aktionsanforderungen

| Aktion              | Mindest-Projektrolle |
| ------------------- | -------------------- |
| Daten lesen         | `viewer`             |
| Datensätze erstellen | `consultant`        |
| Datensätze aktualisieren | `consultant`    |
| Datensätze löschen  | `manager`            |
| Mitglieder verwalten | `manager`           |
| Daten exportieren   | `consultant`         |

## Vertraulichkeits-Flags

Jede Projektberechtigung hat drei unabhängige boolesche Flags, die pro Benutzer und Projekt von einem Manager oder Eigentümer gesetzt werden.

| Flag                  | Steuert                                                              |
| --------------------- | -------------------------------------------------------------------- |
| `can_view_contact`    | Telefonnummern und E-Mail-Adressen in Personendatensätzen            |
| `can_view_personal`   | Vollständiger Name, Geburtsdatum, nationale ID (`external_id`), Einwilligungsinfo |
| `can_view_documents`  | Zugriff auf Dokumentdateien und Dokumentmetadaten                    |

Wenn ein Flag deaktiviert ist, werden die entsprechenden Felder aus API-Antworten und CSV-Exporten ausgelassen — die Daten bleiben in der Datenbank, werden aber nicht an diesen Benutzer gesendet. Ein Berater mit `can_view_personal: false` kann nationale IDs oder Geburtsdaten nicht über einen Export abrufen, auch wenn er Datensätze erstellen kann.

## Berechtigungen zuweisen

Nur Plattformadministratoren und Projekteigentümer/-manager können Projektberechtigungen zuweisen.

### Benutzer zu einem Projekt hinzufügen

```http
POST /admin/projects/:project_id/permissions
```

```json
{
  "user_id": "01J...",
  "role": "consultant",
  "can_view_contact": true,
  "can_view_personal": false,
  "can_view_documents": false
}
```

### Bestehende Berechtigung aktualisieren

```http
PUT /admin/projects/:project_id/permissions/:permission_id
```

### Benutzer aus Projekt entfernen

```http
DELETE /admin/projects/:project_id/permissions/:permission_id
```

## Häufige Einrichtungen

### Feldarbeiter (consultant, eingeschränkte persönliche Daten)

```json
{
  "role": "consultant",
  "can_view_contact": true,
  "can_view_personal": false,
  "can_view_documents": false
}
```

Telefonnummern zur Koordination sichtbar; kein Zugriff auf nationale IDs oder Geburtsdaten.

### Supervisor (manager, Vollzugriff)

```json
{
  "role": "manager",
  "can_view_contact": true,
  "can_view_personal": true,
  "can_view_documents": true
}
```

### Externer Prüfer (viewer, keine sensiblen Daten)

```json
{
  "role": "viewer",
  "can_view_contact": false,
  "can_view_personal": false,
  "can_view_documents": false
}
```
