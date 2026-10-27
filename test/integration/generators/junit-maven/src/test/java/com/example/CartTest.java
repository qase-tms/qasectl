package com.example;

import static org.junit.jupiter.api.Assertions.assertEquals;

import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;

class CartTest {

    @Test
    void addItemPasses() {
        assertEquals(2, 1 + 1);
    }

    @Test
    void totalIsWrongFails() {
        assertEquals(40, 10 + 20, "total must include both items");
    }

    @Disabled("not implemented yet")
    @Test
    void couponIsSkipped() {
        assertEquals(1, 1);
    }
}
