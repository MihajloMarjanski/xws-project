package com.example.agent.service;

import com.example.agent.model.*;
import com.example.agent.repository.AdminRepository;
import com.example.agent.repository.CompanyOwnerRepository;
import com.example.agent.repository.CompanyRepository;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

import javax.servlet.http.HttpServletRequest;
import java.util.*;

@Service
@Slf4j
public class AdminService {

    @Autowired
    private RoleService roleService;
    @Autowired
    private AdminRepository adminRepository;
    @Autowired
    private CompanyRepository companyRepository;
    @Autowired
    private CompanyOwnerRepository companyOwnerRepository;
    @Autowired
    EmailService emailService;


    public ResponseEntity<?> approveCompany(Integer id, HttpServletRequest request) {
        Optional<Company> company = companyRepository.findById(id);
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        if (!company.isPresent()) {
            log.warn("Ip: {}, username: {}, Company doesn't exist!", request.getRemoteAddr(), authentication.getPrincipal().toString());
            return new ResponseEntity<>(HttpStatus.BAD_REQUEST);
        }

        Role role = roleService.findByName("ROLE_COMPANY_OWNER");
        CompanyOwner owner = company.get().getCompanyOwner();
        Set<Role> ownerRoles = new HashSet<>();
        ownerRoles.add(role);
        owner.setRoles(ownerRoles);
        companyOwnerRepository.save(owner);

        company.get().setApproved(true);
        companyRepository.save(company.get());
        log.info("Ip: {}, username: {}, Company created successfully!", request.getRemoteAddr(), authentication.getPrincipal().toString());
        return new ResponseEntity<>(HttpStatus.OK);
    }

    public Collection<String> findAllUsernames() {
        return adminRepository.findAllUsernames();
    }

    public boolean isPinOk(String username, String pin) {
        Admin user = adminRepository.findByUsername(username);
        if (user == null)
            return false;
        Calendar c = Calendar.getInstance();
        c.setTime(user.getPinCreatedDate());
        c.add(Calendar.MINUTE, 1);

        if (user.getPin().equals("") || c.getTime().before(new Date())) {
            return false;
        }
        PasswordEncoder passwordEncoder = new BCryptPasswordEncoder();
        String saltedPin = pin.concat(user.getSalt());
        boolean match = passwordEncoder.matches(saltedPin, user.getPin());
        return match;
    }

    public Admin findByUsername(String username) {
        return adminRepository.findByUsername(username);
    }

    public ResponseEntity<?> getByUsername(String username, HttpServletRequest request) {
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        if (findByUsername(username) == null){
            log.warn("Ip: {}, username: {}, User doesn't exist!", request.getRemoteAddr(), authentication.getPrincipal().toString());
            return new ResponseEntity<>(HttpStatus.BAD_REQUEST);
        }
        return new ResponseEntity<>(findByUsername(username), HttpStatus.OK);
    }

    public ResponseEntity<?> updateAdmin(Admin client, HttpServletRequest request) {
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        Admin admin = findByUsername(client.getUsername());
        admin.setFirstName(client.getFirstName());
        admin.setLastName(client.getLastName());
        if(!admin.getPassword().equals(client.getPassword()) || client.getPassword() == "") {
            PasswordEncoder passwordEncoder = new BCryptPasswordEncoder();
            admin.setPassword(passwordEncoder.encode(client.getPassword().concat(admin.getSalt())));
        }
        adminRepository.save(admin);
        log.info("Ip: {}, username: {}, User updated successfully!", request.getRemoteAddr(), authentication.getPrincipal().toString());
        return new ResponseEntity<>(admin, HttpStatus.OK);
    }

    public void send2factorAuthPin(Admin admin) {
        PasswordEncoder passwordEncoder = new BCryptPasswordEncoder();
        String pin = RandomStringInitializer.generatePin();
        admin.setPin(passwordEncoder.encode(pin.concat(admin.getSalt())));
        admin.setPinCreatedDate(new Date());
        adminRepository.save(admin);
        emailService.send2factorAuthPin(admin.getEmail(), pin);
    }
}
