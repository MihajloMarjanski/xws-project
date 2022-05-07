package com.example.agent.repository;

import com.example.agent.model.CompanyOwner;
import org.springframework.data.jpa.repository.JpaRepository;

public interface CompanyOwnerRepository extends JpaRepository<CompanyOwner, Integer> {
    CompanyOwner findByUsername(String username);
}
