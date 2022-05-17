package com.example.agent.repository;

import com.example.agent.model.CompanyOwner;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;

import java.util.Collection;

public interface CompanyOwnerRepository extends JpaRepository<CompanyOwner, Integer> {
    CompanyOwner findByUsername(String username);

    @Query("SELECT c.username FROM CompanyOwner c")
    Collection<String> findAllUsernames();
}
