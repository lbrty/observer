---
title: Einführung
weight: 1
---

Observer ist eine selbst gehostete Fallverwaltungsplattform für NGOs, die mit Binnenvertriebenen arbeiten.

## Das Problem

Kleine und mittelgroße NGOs, die in der Ukraine mit vertriebenen Personen arbeiten, stehen vor einer Lücke: Papierregister skalieren nicht, aber die Systeme, die es tun — proGres, Primero, ActivityInfo — erfordern UN-Partnerschaften, technische Partner oder SaaS-Abonnements, die für diese Organisationen nicht zugänglich sind.

Observer schließt diese Lücke.

## Was es leistet

Observer gibt Ihrer Organisation ein privates, sicheres System zur Erfassung der Menschen, die Sie betreuen. Es läuft auf Ihrem eigenen Server — kein Cloud-Service, kein Abonnement, kein Dritter sieht jemals Ihre Daten.

Mit Observer kann Ihr Team:

- **Personen und Familien registrieren** — persönliche Daten, Haushaltsbeziehungen und Dokumente erfassen
- **Unterstützung verfolgen** — Beratungen, Überweisungen und die Art der geleisteten Hilfe dokumentieren
- **Bewegungen nachverfolgen** — woher Menschen kamen, wohin sie gezogen sind und warum
- **Zugriff steuern** — festlegen, wer in Ihrem Team was sehen darf, bis hin zu Kontaktdaten und Dokumenten
- **Berichte erstellen** — integrierte Aufschlüsselungen nach Geschlecht, Alter, Region, Unterstützungsart, Vulnerabilitätskategorie und mehr — filterbar nach EU-, USAID- und bilateralen Geberanforderungen

Eine Person mit grundlegenden Serverkenntnissen kann es in unter einer Stunde einrichten.

## Für wen es gedacht ist

Organisationen mit:

- 5 bis 50 Mitarbeitern, einer oder keiner IT-Person
- Einem oder mehreren geberfinanzierten Hilfsprojekten
- Berichtspflichten gegenüber EU, USAID oder bilateralen Gebern
- Datenschutzanforderungen, die die Nutzung von Drittanbieter-SaaS verhindern
- Keiner bestehenden Partnerschaft mit einer UN-Agentur, die Zugang zu Primero oder proGres gewährt

## Warum nicht die Alternativen

| System | Hürde |
| --- | --- |
| **proGres v4** | Erfordert UNHCR-Partnerschaft; nicht selbst hostbar |
| **Primero** | Erfordert UNICEF/IRC als technischen Partner für die Bereitstellung |
| **ActivityInfo** | SaaS mit nutzerbezogener Preisgestaltung; für aggregiertes Monitoring konzipiert, nicht für Fallverwaltung |
| **KoBoToolbox** | Nur Datenerfassung — keine dauerhaften Fallakten |

## Was Observer besonders macht

**Sie brauchen niemandes Erlaubnis.** Kein Partner-Onboarding, keine SaaS-Vereinbarung. Installieren Sie es auf Ihrem Server und beginnen Sie zu arbeiten.

**Ihre Daten bleiben bei Ihnen.** Alles läuft auf Infrastruktur, die Sie kontrollieren. Keine Daten verlassen Ihren Server, es sei denn, Sie exportieren sie.

**Zugriffskontrolle ist integriert.** Plattformrollen (Admin, Mitarbeiter, Berater, Gast) werden mit Projektrollen (Eigentümer, Manager, Berater, Betrachter) und Sensitivitätsstufen kombiniert, die steuern, wer Kontaktdaten, persönliche Details und Dokumente sehen darf.

**Berichte entsprechen dem, was Geber tatsächlich verlangen.** 12 Berichtsdimensionen — Beratungszahlen, IDP-Herkunft, Geschlechts-/Altersdisaggregation, Unterstützungsbereich, Vulnerabilitätskategorie, Region, Büro, Tags, Familieneinheiten und Fallstatus — jeweils filterbar nach Zeitraum, Unterstützungsart und demografischen Kriterien.

**Familien werden als Einheiten erfasst.** Typisierte Haushaltsbeziehungen (Haushaltsvorstand, Ehepartner, Kind, Elternteil, Geschwister) ermöglichen Berichte auf Familienebene.

**IDP-Status wird berechnet, nicht geraten.** Der Vertreibungsstatus einer Person wird aus ihrem Herkunftsort und der Frage, ob dieses Gebiet eine Konfliktzone ist, abgeleitet — keine manuelle Klassifizierung nötig.

## Unterstützte Sprachen

Die Benutzeroberfläche wird mit sechs Sprachen ausgeliefert: Englisch, Ukrainisch, Russisch, Deutsch, Türkisch und Kirgisisch (lateinische Schrift). Kirgisisch verwendet eine eigene lateinische Transliteration, da das offizielle kirgisische Lateinalphabet 2023 eingeführt wurde und gängige Übersetzungstools es noch nicht unterstützen — wir pflegen eigene Transliterationsregeln, um akkuraten, natürlich klingenden Text für zentralasiatische Einsätze bereitzustellen.

## Was Observer nicht leistet

Observer ist nicht konzipiert für:

- Interinstitutionelle Überweisungs-Workflows (Primero-Bereich)
- Biometrische Deduplizierung
- Cluster-Level 5W/3W-Aggregatberichte (ActivityInfo-Bereich)

Wenn Sie Cluster-Dashboards benötigen, exportieren Sie Ihre Daten von Observer nach ActivityInfo.
