package br.com.condosnap.buzufba.repository;

import java.util.Optional;

import org.springframework.data.jpa.repository.JpaRepository;

import br.com.condosnap.buzufba.entity.Location;

public interface LocationRepository extends JpaRepository<Location, Long>{
    Optional<Location> findByName(String name);
}
