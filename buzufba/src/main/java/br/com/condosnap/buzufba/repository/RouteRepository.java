package br.com.condosnap.buzufba.repository;

import java.util.List;

import org.springframework.data.jpa.repository.JpaRepository;

import br.com.condosnap.buzufba.entity.Route;
import br.com.condosnap.buzufba.projection.RouteNames;

public interface RouteRepository extends JpaRepository<Route, String>{
    List<Route> findAll();

    List<RouteNames> findAllBy();
}
