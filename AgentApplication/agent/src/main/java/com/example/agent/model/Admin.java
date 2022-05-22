package com.example.agent.model;

import com.fasterxml.jackson.annotation.JsonIgnore;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.userdetails.UserDetails;

import javax.persistence.*;
import java.util.*;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Entity
public class Admin extends User implements UserDetails {

//    @ManyToOne(fetch = FetchType.EAGER)
//    @JoinColumn(name = "role_id")
//    private Role role;

    @ManyToMany(fetch = FetchType.EAGER)
    @JoinTable(name = "admin_roles",
            joinColumns = @JoinColumn(name = "admin_id", referencedColumnName = "id"),
            inverseJoinColumns = @JoinColumn(name = "role_id", referencedColumnName = "id"))
    private Set<Role> roles = new HashSet<>();


    @JsonIgnore
    @Override
    public Collection<? extends GrantedAuthority> getAuthorities() {
        List<GrantedAuthority> authorities = new ArrayList<>();
        for(Role r : this.roles) {
            authorities.add(new SimpleGrantedAuthority(r.name));
            for(Permission p : r.getPermissions())
                authorities.add(new SimpleGrantedAuthority(p.name));
        }

        return authorities;
    }

    @JsonIgnore
    @Override
    public boolean isAccountNonExpired() {
        return true;
    }

    @JsonIgnore
    @Override
    public boolean isAccountNonLocked() {
        return true;
    }

    @JsonIgnore
    @Override
    public boolean isCredentialsNonExpired() {
        return true;
    }

    @JsonIgnore
    @Override
    public boolean isEnabled() {
        return this.isActivated;
    }
}
