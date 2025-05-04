package br.com.condosnap.buzufba.controller;

import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import br.com.condosnap.buzufba.entity.Route;
import br.com.condosnap.buzufba.projection.RouteNames;
import br.com.condosnap.buzufba.service.RouteService;

import java.util.List;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestParam;


@RestController
@RequestMapping("/route")
public class RouteController {
    
    @Autowired
    private RouteService routeService;


    @GetMapping("/all")
    public List<Route> getAllAvailableRoutes() {
        return routeService.getAllRoutes();
    }


    @GetMapping("/all/names")
    public List<RouteNames> getAllAvailableRouteNames() {
        return routeService.getAllAvailableRouteNames();
    }
    

    @GetMapping("/{id}")
    public Route getMethodName(@PathVariable String id) {
        return routeService.getRouteById(id);
    }
    
    
}
