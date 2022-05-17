package com.example.agent.repository;

import com.example.agent.model.Admin;
import com.example.agent.model.Client;
import org.springframework.data.jpa.repository.JpaRepository;

public interface ClientRepository extends JpaRepository<Client, Integer> {
}
