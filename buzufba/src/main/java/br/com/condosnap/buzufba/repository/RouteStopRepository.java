package br.com.condosnap.buzufba.repository;

import org.springframework.data.jpa.repository.JpaRepository;

import br.com.condosnap.buzufba.entity.RouteStop;

public interface RouteStopRepository extends JpaRepository<RouteStop, Long> {
}