DO $$
BEGIN
    FOR i IN 1..10000 LOOP
        INSERT INTO wallet_items (wallet_id, symbol, quantity) VALUES
            ('wallet' || i, 'BTCUSD', random()),
            ('wallet' || i, 'ETHUSD', random());
    END LOOP;
END $$;
