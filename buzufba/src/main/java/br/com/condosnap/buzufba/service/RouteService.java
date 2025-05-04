package br.com.condosnap.buzufba.service;

import java.util.List;
import java.util.Map;
import java.util.Optional;

import org.apache.catalina.connector.Response;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import br.com.condosnap.buzufba.entity.Route;
import br.com.condosnap.buzufba.projection.RouteNames;
import br.com.condosnap.buzufba.repository.RouteRepository;

@Service
public class RouteService {
    
    @Autowired
    private RouteRepository routeRepository;

    public List<Route> getAllRoutes(){
        return routeRepository.findAll();
    }

    public List<RouteNames> getAllAvailableRouteNames(){
        return routeRepository.findAllBy();
    }   


    public Route getRouteById(String id){
        Optional<Route> route = routeRepository.findById(id);
        
        if (route.isEmpty()) {
            throw new RuntimeException("Rota não disponível!");
        }
        
        return route.get();
    }
}
