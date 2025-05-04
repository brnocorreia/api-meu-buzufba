package br.com.condosnap.buzufba.entity;

import java.util.ArrayList;
import java.util.List;

import com.fasterxml.jackson.annotation.JsonManagedReference;

import jakarta.persistence.CascadeType;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.ManyToOne;
import jakarta.persistence.OneToMany;
import jakarta.persistence.Table;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Setter
@Getter
@NoArgsConstructor
@AllArgsConstructor
@Entity
@Table(name = "route")
public class Route {

    @Id
    @Column(length = 50)
    private String id;

    @Column(nullable = false)
    private String name;

    private double tripLength;

    @ManyToOne
    @JoinColumn(name = "departure_location_id")
    private Location departureLocation;

    @ManyToOne
    @JoinColumn(name = "arrival_location_id")
    private Location arrivalLocation;

    @Column(length = 1024)
    private List<String> notes;

    @OneToMany(mappedBy = "route", cascade = CascadeType.ALL, orphanRemoval = true)
    @JsonManagedReference
    private List<RouteStop> stops = new ArrayList<>();

    @OneToMany(mappedBy = "route", cascade = CascadeType.ALL, orphanRemoval = true)
    @JsonManagedReference
    private List<DepartureTime> departures = new ArrayList<>();
}
