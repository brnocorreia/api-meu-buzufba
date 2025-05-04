package br.com.condosnap.buzufba.repository;

import org.springframework.data.jpa.repository.JpaRepository;

import br.com.condosnap.buzufba.entity.DepartureTime;

public interface DepartureTimeRepository extends JpaRepository<DepartureTime, Long>{
    
}
