package com.example.agent.repository;

import com.example.agent.model.Admin;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;

import java.util.Collection;

public interface AdminRepository extends JpaRepository<Admin, Integer> {
    Admin findByUsername(String username);

    @Query("SELECT c.username FROM Admin c")
    Collection<String> findAllUsernames();
}
