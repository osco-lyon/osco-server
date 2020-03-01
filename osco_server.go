/*
 * Authors : Thomas Andre, Victor Bonnin, Pierre Niogret, Bénédicte Thomas
 * License: AGPLv3 or later
 */

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

/*
 * Structures de JSON personnalisées
 */

type jsonStructPerso struct {
	GEOMETRY []geometry
	CRITERE  critere
}

type critere struct {
	SURFACE     int
	AMENAGEMENT int
	VITESSE     int
}

type contribution struct {
	TAGS []string
	BBOX []float64
	NELT int
}

type geometry struct {
	X float64
	Y float64
}

/*
 * Programme principal
 */
func main() {
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("/var/www/osco_html"))))
	http.HandleFunc("/api/perso", handleAPIperso)
	http.HandleFunc("/api/rapide", handleAPIrapide)
	http.HandleFunc("/api/recommande", handleAPIrecommande)
	http.HandleFunc("/api/amenagement", handleAPIamenagement)
	http.HandleFunc("/api/contribution", handleContribution)
	log.Println("Listening...")
	//http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/osco.anatidaepho.be/fullchain.pem", "/etc/letsencrypt/live/osco.anatidaepho.be/privkey.pem", nil)
	http.ListenAndServe(":80", nil)
}

/*
 * Fonction de gestion /api/recommande
 */
func handleAPIrecommande(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	r.ParseForm()
	geometryJson := r.Form.Get("json")
	var geometries []jsonStructPerso
	json.Unmarshal([]byte(geometryJson), &geometries)
	x1 := geometries[0].GEOMETRY[0].X
	y1 := geometries[0].GEOMETRY[0].Y
	x2 := geometries[0].GEOMETRY[1].X
	y2 := geometries[0].GEOMETRY[1].Y
	reqIti := fmt.Sprintf("SELECT itineraire_secu(%f,%f,%f,%f)\n", x1, y1, x2, y2)
	recommande := pgQuery(reqIti)
	fmt.Fprintf(w, recommande)
}

/*
 * Fonction de gestion /api/amenagement
 */
func handleAPIamenagement(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	r.ParseForm()
	geometryJson := r.Form.Get("json")
	var geometries []jsonStructPerso
	json.Unmarshal([]byte(geometryJson), &geometries)
	x1 := geometries[0].GEOMETRY[0].X
	y1 := geometries[0].GEOMETRY[0].Y
	x2 := geometries[0].GEOMETRY[1].X
	y2 := geometries[0].GEOMETRY[1].Y
	reqIti := fmt.Sprintf("SELECT itineraire_amenag(%f,%f,%f,%f)\n", x1, y1, x2, y2)
	amenagement := pgQuery(reqIti)
	fmt.Fprintf(w, amenagement)
}

/*
 * Fonction de gestion /api/rapide
 */
func handleAPIrapide(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	r.ParseForm()
	geometryJson := r.Form.Get("json")
	var geometries []jsonStructPerso
	json.Unmarshal([]byte(geometryJson), &geometries)
	x1 := geometries[0].GEOMETRY[0].X
	y1 := geometries[0].GEOMETRY[0].Y
	x2 := geometries[0].GEOMETRY[1].X
	y2 := geometries[0].GEOMETRY[1].Y
	reqIti := fmt.Sprintf("SELECT itineraire_court(%f,%f,%f,%f)\n", x1, y1, x2, y2)
	rapide := pgQuery(reqIti)
	fmt.Fprintf(w, rapide)
}

/*
 * Fonction de gestion /api/perso
 */
func handleAPIperso(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	r.ParseForm()
	geometryJson := r.Form.Get("json")
	fmt.Printf(geometryJson)
	var geometries []jsonStructPerso
	json.Unmarshal([]byte(geometryJson), &geometries)
	x1 := geometries[0].GEOMETRY[0].X
	y1 := geometries[0].GEOMETRY[0].Y
	x2 := geometries[0].GEOMETRY[1].X
	y2 := geometries[0].GEOMETRY[1].Y
	qualSurface := geometries[0].CRITERE.SURFACE
	vitesseVehi := geometries[0].CRITERE.VITESSE
	qualCycli := geometries[0].CRITERE.AMENAGEMENT
	reqIti := fmt.Sprintf("SELECT itineraire_perso(%f,%f,%f,%f,%v,%v,%v)\n", x1, y1, x2, y2, qualSurface, vitesseVehi, qualCycli)
	perso := pgQuery(reqIti)
	fmt.Fprintf(w, perso)
}

/*
 * Fonction de gestion /api/contribution
 */
func handleContribution(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	r.ParseForm()
	contributionJson := r.Form.Get("json")
	var contributions []contribution
	json.Unmarshal([]byte(contributionJson), &contributions)
	x1 := contributions[0].BBOX[0]
	y1 := contributions[0].BBOX[1]
	x2 := contributions[0].BBOX[2]
	y2 := contributions[0].BBOX[3]
	tags := contributions[0].TAGS
	taglist := "{"
	for _, tag := range tags {
		taglist += tag
		taglist += ","
	}
	taglist = taglist[:len(taglist)-1] + "}"
	nelt := contributions[0].NELT
	reqContribution := fmt.Sprintf("SELECT contrib_items(%f,%f,%f,%f,'%s'::varchar[],%d)\n", x1, y1, x2, y2, taglist, nelt)
	contribution := pgQuery(reqContribution)
	fmt.Fprintf(w, contribution)
}

/*
 * Fonction de requête vers la base PostgreSQL
 */
func pgQuery(query string) string {
	connStr := "user=osco dbname=osco password=20GeoNum20"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	defer db.Close()
	var itineraire string
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&itineraire)
		if err != nil {
			log.Fatal(err)
			return ""
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return itineraire
}

/*
 * Fonction ajoutant l'en-tête CORS nécessaire pour le requêtage distant
 */
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
