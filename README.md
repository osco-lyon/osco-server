# Serveur OSCO

Partie serveur de l'application OSCO, développée en Go. Celui-ci fait appel, par défaut, à une base de données PostgreSQL osco hébergée sur la même machine.

## Dépendences

Il est nécessaire, pour lancer l'application, de disposer du connecteur PostgreSQL [pq](https://github.com/lib/pq). Celui-ci peut être installé par la commande suivante :

```bash
go get github.com/lib/pq
```

## Configuration

Le [client html](https://github.com/osco-lyon/osco-client) doit être placé dans le dossier /var/www/osco_html.

La base de données utilisée repose sur les [scripts d'installation du serveur](https://github.com/osco-lyon/osco-install).
