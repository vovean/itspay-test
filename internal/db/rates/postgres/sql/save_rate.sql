insert into rates (ask, bid, received_at)
values ($1, $2, $3)
on conflict (received_at) do update set ask=excluded.ask,
                                        bid=excluded.bid