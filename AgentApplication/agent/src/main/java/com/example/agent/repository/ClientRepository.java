package com.example.agent.repository;

import com.example.agent.model.Admin;
import com.example.agent.model.Client;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;

import java.util.Collection;

public interface ClientRepository extends JpaRepository<Client, Integer> {
    Client findByEmail(String email);

    Client findByUsername(String username);

    @Query("SELECT c.username FROM Client c")
    Collection<String> findAllUsernames();
}
