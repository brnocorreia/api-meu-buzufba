package br.com.condosnap.buzufba.seeder;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.Set;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.CommandLineRunner;
import org.springframework.stereotype.Component;

import br.com.condosnap.buzufba.entity.DepartureTime;
import br.com.condosnap.buzufba.entity.Location;
import br.com.condosnap.buzufba.entity.Route;
import br.com.condosnap.buzufba.entity.RouteStop;
import br.com.condosnap.buzufba.repository.DepartureTimeRepository;
import br.com.condosnap.buzufba.repository.LocationRepository;
import br.com.condosnap.buzufba.repository.RouteRepository;
import br.com.condosnap.buzufba.repository.RouteStopRepository;

@Component
public class DatabaseSeeder implements CommandLineRunner {

    @Autowired
    private LocationRepository locationRepository;
    @Autowired
    private RouteRepository routeRepository;
    @Autowired
    private RouteStopRepository routeStopRepository;
    @Autowired
    private DepartureTimeRepository departureTimeRepository;

    @Override
    public void run(String... args) {
        // Locais da rota
        List<String> locationNames = List.of(
                "Estacionamento PAF I - Matemática", "Av. Garibaldi", "Campus Vale do Canela",
                "Viaduto Campo Grande", "Avenida 7 de Setembro - Faculdade de Economia",
                "Belas Artes", "Reitoria", "Creche – Canela", "Politécnica",
                "Arquitetura", "Instituto de Geociências", "Circular",
                "São Lázaro", "Viaduto Federação", "Residência 5",
                "Ondina/PAF1", "Residência Universitária Garibaldi", "Deli&Cia", "Direito",
                "Música", "ISC", "Odontologia", "Nutrição", "Geociências", "Piedade", "Centro de Esportes","Portaria Principal", "Proae","Facom", "Reitoria sentido Campo Grande","Retorno - Rua Baronesa de Sauípe", "Av. Garibaldi - Ponto R5");

        // Obtém ou cria os locais
        Map<String, Location> locationMap = new HashMap<>();
        for (String name : locationNames) {
            Location location = locationRepository.findByName(name)
                    .orElseGet(() -> locationRepository.save(new Location(name)));
            locationMap.put(name, location);
        }

        seedRoute("EXPRESSO", "Rota Expresso", 13,
                List.of("6h30", "7h30", "8h40", "9h50", "11h00", "12h20",
                        "13h20", "14h30", "15h40", "16h50", "18h00",
                        "19h10", "20h10", "21h10", "22h30"),
                locationMap.get("Estacionamento PAF I - Matemática"),
                locationMap.get("Circular"),
                List.of("20h10 é o último horário a entrar na Piedade"),
                List.of("Estacionamento PAF I - Matemática", "Av. Garibaldi", "Campus Vale do Canela",
                        "Viaduto Campo Grande", "Avenida 7 de Setembro - Faculdade de Economia",
                        "Belas Artes", "Reitoria", "Creche – Canela", "Politécnica",
                        "Arquitetura", "Instituto de Geociências"),
                Set.of("Estacionamento PAF I - Matemática", "Av. Garibaldi", "Campus Vale do Canela",
                        "Viaduto Campo Grande", "Avenida 7 de Setembro - Faculdade de Economia", "Belas Artes"),
                Set.of("Reitoria", "Creche – Canela", "Politécnica",
                        "Arquitetura", "Instituto de Geociências", "Estacionamento PAF I - Matemática"));

        seedRoute("B1", "Rota B1", 11,
                List.of("6h10", "7h00", "8h00", "9h00", "10h00", "11h00",
                        "12h00", "13h00", "15h00", "16h00", "17h00",
                        "18h00", "19h00", "20h30", "21h40", "22h20"),
                locationMap.get("São Lázaro"),
                locationMap.get("Reitoria"),
                List.of("Após fechamento de São Lázaro, carro volta para Ondina e retoma rota até último horário."),
                List.of("São Lázaro", "Politécnica", "Arquitetura", "Viaduto Federação",
                        "Residência 5", "Instituto de Geociências", "Estacionamento PAF I - Matemática",
                        "Av. Garibaldi", "Campus Vale do Canela Entrada ICS", "Viaduto Campo Grande",
                        "Belas Artes", "Reitoria", "Creche – Canela"),
                Set.of("São Lázaro", "Politécnica", "Arquitetura", "Viaduto Federação",
                        "Residência 5", "Instituto de Geociências", "Estacionamento PAF I - Matemática"),
                Set.of("Av. Garibaldi", "Campus Vale do Canela Entrada ICS", "Viaduto Campo Grande",
                        "Belas Artes", "Reitoria", "Creche – Canela", "Politécnica", "São Lázaro"));

        seedRoute("B2", "Rota B2", 13,
                List.of("6h00", "7h00", "8h00", "9h00", "10h00", "11h00",
                        "12h00", "13:30", "14h30", "16h00", "17h40",
                        "18h30", "19h50", "20h30", "21h40", "22h30"),
                locationMap.get("Ondina/PAF1"),
                locationMap.get("Reitoria"),
                List.of("19h50 é o último horário a entrar em São Lázaro"),
                List.of("Ondina/PAF1", "Residência Universitária Garibaldi", "Residência I - Vitória",
                        "Deli&Cia", "Reitoria", "Creche – Canela", "Politécnica", "Arquitetura",
                        "Instituto de Geociências"),
                Set.of("Ondina/PAF1", "Residência Universitária Garibaldi", "Arquitetura",
                        "São Lázaro", "Politécnica", "Creche – Canela", "Reitoria"),
                Set.of("Residência I - Vitória", "Deli&Cia", "Politécnica", "São Lázaro",
                        "Arquitetura", "Instituto de Geociências", "Ondina/PAF1"));

        seedRoute("B3", "Rota B3", 15.5,
                List.of("6h30", "7h10", "8h40", "9h50", "11h00", "12h10",
                        "13h20", "14h30", "15h40", "16h50", "18h00",
                        "19h10", "20h30", "21h20", "22h20"),
                locationMap.get("Direito"),
                locationMap.get("Ondina/PAF1"),
                List.of("19h10 é o último horário a entrar em São Lázaro"),
                List.of("Direito", "Música", "ISC", "Odontologia", "Nutrição",
                        "Ondina/PAF1", "Residência Universitária Garibaldi", "Deli&Cia",
                        "Reitoria", "Creche – Canela", "Politécnica", "Arquitetura", "Instituto de Geociências"),
                Set.of("Estacionamento PAF I - Matemática", "Av. Garibaldi - Ponto R5", "Arquitetura",
                        "São Lázaro", "Politécnica", "Creche – Canela", "Reitoria sentido Campo Grande",
                        "Retorno - Rua Baronesa de Sauípe", "Belas Artes", "Reitoria",
                        "Deli&Cia", "Direito"),
                Set.of("Escola de Música - ISC - Odontologia - Nutrição", "Reitoria",
                        "Politécnica", "Arquitetura", "Instituto de Geociências",
                        "Estacionamento PAF I - Matemática"));

        seedRoute("B4", "Rota B4", 14,
                List.of("6h20", "7h20", "8h20", "9h30", "10h40", "11h40",
                        "12h40", "14h00", "15h20", "16h30", "17h40",
                        "20h00", "21h20", "22h30"),
                locationMap.get("Ondina/PAF1"),
                locationMap.get("Piedade"),
                List.of("18:50 é o último horário a entrar em São Lázaro"),
                List.of("Ondina/PAF1", "Residência Universitária Garibaldi", "Reitoria",
                        "Economia", "Belas Artes", "São Lázaro", "Creche – Canela",
                        "Politécnica", "Arquitetura", "Instituto de Geociências"),
                Set.of("Estacionamento PAF I - Matemática", "Av. Garibaldi - Ponto R5",
                        "Arquitetura", "Politécnica", "Creche – Canela", "Reitoria", "Piedade"),
                Set.of("Piedade", "Belas Artes", "Reitoria", "Creche – Canela",
                        "Politécnica", "São Lázaro", "Arquitetura", "Instituto de Geociências",
                        "Estacionamento PAF I - Matemática"));

        seedRoute("B5", "Rota B5", 17,
                List.of("6h20", "7h20", "8h40", "10h00", "11h20", "12h40",
                        "14h00", "15h20", "16h40", "18h00", "19h20",
                        "20h40", "22h20"),
                locationMap.get("Facom"),
                locationMap.get("Reitoria"),
                List.of("19h20 é o último horário a entrar em São Lázaro"),
                List.of("Facom", "Residência Universitária Garibaldi", "Arquitetura", "São Lázaro",
                        "Politécnica", "Creche – Canela", "Reitoria", "Instituto de Geociências",
                        "Residência I - Vitória", "Deli&Cia", "Ondina/PAF1"),
                Set.of("Instituto de Geociências", "Facom", "Portaria Principal", "Centro de Esportes",
                        "Av. Garibaldi - Ponto R5", "Proae", "São Lázaro", "Politécnica",
                        "Creche – Canela", "Reitoria"),
                Set.of("Campo Grande", "Residência I - Ponto de Distribuição Vitória",
                        "Deli&Cia - acesso direito", "Politécnica", "São Lázaro", "Arquitetura",
                        "Estacionamento PAF I - Matemática", "Facom", "Instituto de Geociências"));

    }

    private void seedRoute(String id, String name, double tripLength,
            List<String> departures, Location departureLocation,
            Location arrivalLocation, List<String> notes,
            List<String> servedLocations, Set<String> departureStops,
            Set<String> arrivalStops) {
        // Criação da rota
        Route route = new Route();
        route.setId(id);
        route.setName(name);
        route.setTripLength(tripLength);
        route.setDepartureLocation(departureLocation);
        route.setArrivalLocation(arrivalLocation);
        route.setNotes(notes);
        route = routeRepository.save(route);

        final Route savedRoute = route;

        departures.forEach(time -> {
            DepartureTime departure = new DepartureTime();
            departure.setRoute(savedRoute);
            departure.setTime(time);
            departureTimeRepository.save(departure);
        });

        List<String> orderedStops = new ArrayList<>(departureStops);
        int order = 0;
        for (String stopName : orderedStops) {
            Optional<Location> locationOpt = locationRepository.findByName(stopName);
            if (locationOpt.isEmpty()) {
                throw new IllegalStateException("Local não encontrado: " + stopName);
            }

            RouteStop stop = new RouteStop();
            stop.setRoute(route);
            stop.setLocation(locationOpt.get());
            stop.setStopOrder(order++);
            stop.setDeparture(departureStops.contains(stopName));
            stop.setArrival(arrivalStops.contains(stopName));
            routeStopRepository.save(stop);
        }
    }
}